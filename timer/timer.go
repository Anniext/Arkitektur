package timer

import (
	"container/list"
	"github.com/Anniext/Arkitektur/system/log"
	"math"
	"runtime"
	"time"
)

const EBucketBits = 6
const EBucketSize = 1 << EBucketBits
const ETimerWheelDepth = 3

//const EBucketMask = EBucketSize - 1
//const ETimerWheelSize = EBucketSize * ETimerWheelDepth

const EDefaultTickMs = 1000

const ETimingWheelStatusDead = 0
const ETimingWheelStatusRun = 1

var EMaxConcurrentTaskNum = runtime.NumCPU()

/*
最小精度<毫秒>  Ms代表毫秒
固定3层 64格 如需调整重构创建即可
Element.Value = &TaskData{}
TaskData.element = &Element{}
bucket <- TaskData
overflow <- TaskData

例子：
1	0 1 2 3 4 5 6 7
2	0 1 2 3 4 5 6 7

tw       0,0    1,0    7,0    0,1    0,7    3,4    7,7
task  1  1,_    2,_    _,1-0  1,_    1,_    4,_    _,0-0
      7  7,_    _.1-0  _,1-6  7,_    7,_    _,5-2  _,0-6
      8  _,1-0  _,1-1  _,1-7  _,2-0  _,0-0  _,5-3  _,0-7
      40 _,5-0  _,5-1  _,5-7  _,6-0  _,4-0  _,1-3  _,4-7
      63 _,7-7  _,x-0  _,x-6  _,x-7  _,x-55 _,x-34 _,x-62
      64 _,x-0  _,x-1  _,x-7  _,x-8  _,x-56 _,x-35 _,x-63

*/

// Task 任务参数
type Task struct {
	ExpiresMs     int64  // 过期时间 毫秒
	ExpiresLoopMs int64  // 循环时间 毫秒
	Callback      func() // 回调
}

// gTimingWheel 默认定时器
var gTimingWheel *TimingWheel

// TimingWheel 时间轮
type TimingWheel struct {
	tickMs             int64                                     // 最小间隔
	status             int32                                     // 运行状态
	slot               [ETimerWheelDepth]int64                   // 当前槽位
	bucket             [ETimerWheelDepth][EBucketSize]*list.List // 桶
	overflow           *list.List                                // 溢出桶
	deleteList         *list.List                                // 删除桶
	setCurrentSlotChan chan int64                                // 设置当前槽数
	addTaskChan        chan *TaskData                            // 增加任务
	delTaskChan        chan *TaskData                            // 删除任务
	exitChan           chan struct{}                             // 时间轮退出信号
	waitExitChan       chan struct{}                             // 等待时间轮退出信号
	ticker             *time.Ticker                              // 时间轮转动频率
	submit             func(func()) error                        // 执行任务
	currentSlot        int64                                     // 当前槽数
	nextExpiry         int64                                     // 后期优化空转使用
}

// SetCurrentSlot 启动设置时强行设置当前时间槽数，调用者处理精度误差, 用于启动时注册大量任务，不明白用途不要调用。   <0 需要延迟毫秒数，  >=0 恢复到当前时间
func (t *TimingWheel) SetCurrentSlot(slot int64) {
	t.setCurrentSlotChan <- slot
}

// AddTimer 创建一个任务
func (t *TimingWheel) AddTimer(task *Task) *TaskData {
	if task == nil {
		return nil
	}

	taskData := &TaskData{
		Element:         nil,
		List:            nil,
		ExpiresSlot:     task.ExpiresMs / t.tickMs,
		ExpiresLoopSlot: task.ExpiresLoopMs / t.tickMs,
		Callback:        task.Callback,
	}
	// 无法处理
	if task.ExpiresLoopMs != 0 && task.ExpiresLoopMs < t.tickMs {
		return nil
	}

	// 小于精度立即执行
	if task.ExpiresMs < t.tickMs {
		err := t.submit(task.Callback)
		if err != nil {
			return nil
		} else if task.ExpiresLoopMs != 0 {
			taskData.ExpiresSlot = taskData.ExpiresLoopSlot
			t.addTaskChan <- taskData
			return taskData
		} else {
			return taskData
		}
	}
	t.addTaskChan <- taskData
	return taskData
}

// DelTimer 删除一个任务
func (t *TimingWheel) DelTimer(taskData *TaskData) {
	if taskData == nil {
		return
	}
	t.delTaskChan <- taskData
}

// ResetTimer 重置一个任务
func (t *TimingWheel) ResetTimer(taskData *TaskData, expiresMs, expiresLoopMs int64) *TaskData {
	if taskData == nil {
		return nil
	}
	if expiresLoopMs != 0 && expiresLoopMs < t.tickMs {
		return nil
	}

	t.delTimer(taskData)

	taskData.ExpiresSlot = expiresMs / t.tickMs
	taskData.ExpiresLoopSlot = expiresLoopMs / t.tickMs

	if expiresMs < t.tickMs {
		err := t.submit(taskData.Callback)
		if err != nil {
			return nil
		} else if expiresLoopMs != 0 {
			taskData.ExpiresSlot = taskData.ExpiresLoopSlot
			t.addTaskChan <- taskData
			return taskData
		} else {
			return taskData
		}
	}
	t.addTaskChan <- taskData
	return taskData
}

// Stop 停止时间轮
func (t *TimingWheel) Stop() {
	t.exitChan <- struct{}{}
	_, ok := <-t.waitExitChan
	if !ok {
		log.Infoln("TimingWheel stop error")
	}
}

// Run 启动定时器
func (t *TimingWheel) Run() {
	go (func() {
		t.setCurrentSlot(0)
		t.ticker = time.NewTicker(time.Duration(t.tickMs) * time.Millisecond)
		t.status = ETimingWheelStatusRun
		defer func() {
			t.ticker.Stop()
			t.ticker = nil
			t.waitExitChan <- struct{}{}
		}()
		for {
			select {
			case <-t.ticker.C:
				t.tick(0)
			case v, ok := <-t.setCurrentSlotChan:
				if ok {
					t.setCurrentSlot(v)
				}
			case v, ok := <-t.addTaskChan:
				//log.Println("add chan", v.expiresSlot, v)
				if ok {
					t.addTimer(v)
				}
			case v, ok := <-t.delTaskChan:
				//log.Println("del chan", v.expiresSlot, v)
				if ok {
					t.delTimer(v)
				}
			case <-t.exitChan:
				if t.status == ETimingWheelStatusRun {
					t.clearTimer()
					t.status = ETimingWheelStatusDead
				}
				return
			}
		}
	})()
}

// addTimer 内部添加定时任务
func (t *TimingWheel) addTimer(taskData *TaskData) {
	expiresSlot := taskData.ExpiresSlot
	var currentExpiresSlot, quotient, remainder int64
	var subExpiresSlot float64

	currentExpiresSlot = expiresSlot

	for depth := 0; depth < ETimerWheelDepth; depth++ {
		slot := t.slot[depth]

		if depth == ETimerWheelDepth-1 {
			// 未溢出
			if currentExpiresSlot < EBucketSize {
				currentExpiresSlot = (currentExpiresSlot + slot) % EBucketSize
				taskData.Element = t.bucket[depth][currentExpiresSlot].PushBack(taskData)
				taskData.List = t.bucket[depth][currentExpiresSlot]
				taskData.ExpiresSlot = int64(subExpiresSlot)
				//log.Println("addTimer ", depth, currentExpiresSlot, taskData.expiresSlot, taskData)
				return
			} else {
				taskData.Element = t.overflow.PushBack(taskData)
				taskData.List = t.overflow
				taskData.ExpiresSlot = int64(subExpiresSlot) + int64(float64(currentExpiresSlot-(EBucketSize-slot))*math.Pow(float64(EBucketSize), float64(depth)))
				//log.Println("addTimer ", -1, -1, taskData.expiresSlot, taskData)
				return
			}
		} else {
			currentExpiresSlot += slot
			// 当前层足够
			if currentExpiresSlot < EBucketSize {
				taskData.Element = t.bucket[depth][currentExpiresSlot].PushBack(taskData)
				taskData.List = t.bucket[depth][currentExpiresSlot]
				taskData.ExpiresSlot = int64(subExpiresSlot)
				//log.Println("addTimer ", depth, currentExpiresSlot, taskData.expiresSlot, taskData)
				return
			} else {
				quotient = currentExpiresSlot / EBucketSize
				remainder = currentExpiresSlot % EBucketSize
				subExpiresSlot += float64(remainder) * math.Pow(float64(EBucketSize), float64(depth))
				currentExpiresSlot = quotient
			}
		}
	}
}

// delTimer 内部删除定时任务
func (t *TimingWheel) delTimer(taskData *TaskData) {
	e := taskData.Element
	if e != nil && taskData.List != nil {
		taskData.List.Remove(e)
		//log.Println("delTimer", taskData)
		taskData.Element = nil
		taskData.List = nil
	}
}

// moveTimer 移动任务 用于降轮
func (t *TimingWheel) moveTimer(l *list.List) {
	// 边删除，边添加  防止死循环，或者报错
	var next *list.Element
	if l != nil {
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			taskData := l.Remove(e).(*TaskData)
			taskData.List = nil
			taskData.Element = nil
			t.addTimer(taskData)
		}
	}
}

// setCurrentSlot 设置当前槽数
func (t *TimingWheel) setCurrentSlot(slot int64) {
	if slot < 0 {
		t.currentSlot = slot
	} else {
		cTime := time.Now()
		cTimeNs := cTime.UnixNano()
		cTimeMs := timeToMs(cTime) + 1
		sleepTime := (t.tickMs-(cTimeMs%t.tickMs))*1e6 + (1e6 - (cTimeNs % 1e6))
		time.Sleep(time.Duration(sleepTime))
		t.currentSlot = (cTimeMs / t.tickMs) + 1
		cTimeMsNew := time.Now().UnixNano()
		log.Infoln("setCurrentSlot 0 adjust", t.currentSlot, sleepTime+cTimeNs, cTimeMsNew)
	}
}

// clearTimer 清空所有任务
func (t *TimingWheel) clearTaskList(l *list.List) {
	var next *list.Element
	if l != nil {
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			taskData := l.Remove(e).(*TaskData)
			// 清理task
			taskData.List = nil
			taskData.Element = nil
		}
	}
}

// clearTimer 清空所有任务
func (t *TimingWheel) clearTimer() {
	for depth := 0; depth < ETimerWheelDepth; depth++ {
		for slot := 0; slot < EBucketSize; slot++ {
			l := t.bucket[depth][slot]
			t.clearTaskList(l)
		}
	}
	t.clearTaskList(t.overflow)
}

// tick tick
func (t *TimingWheel) tick(depth int) {
	if depth == 0 {
		// todo t.currentSlot 需要考虑cpu负载导致slot偏移问题并自动修复
		t.currentSlot += 1
		//global.Sugar.Infoln("tick ", t.currentSlot, timer.Now().Unix())
		// 延迟启动 用于启动注册大量任务
		if t.currentSlot <= 0 {
			return
		}
	}

	t.slot[depth] += 1
	// 本层结束 开始降轮
	if t.slot[depth] >= EBucketSize {
		t.slot[depth] = 0

		if depth+1 == ETimerWheelDepth {
			// 处理溢出轮
			t.deleteList, t.overflow = t.overflow, t.deleteList
			t.moveTimer(t.deleteList)
			return
		} else {
			t.tick(depth + 1)
			// 降轮
			l := t.bucket[depth+1][t.slot[depth+1]]

			t.moveTimer(l)
		}

	} else {
	}

	if depth == 0 {
		var next *list.Element
		slot := t.slot[depth]
		t.deleteList, t.bucket[depth][slot] = t.bucket[depth][slot], t.deleteList

		l := t.deleteList
		// 执行函数，删除链表
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			taskData := l.Remove(e).(*TaskData)
			taskData.List = nil
			taskData.Element = nil
			err := t.submit(taskData.Callback)
			if err != nil {
				log.Infoln("submit error", err.Error())
				return
			}
			if taskData.ExpiresLoopSlot > 0 {
				taskData.ExpiresSlot = taskData.ExpiresLoopSlot
				t.addTimer(taskData)
			}
		}
	}
}

// NewTimingWheel 新建: 一个进程不建议使用多个定时器， 除非精度相差很大
func NewTimingWheel(tickMs int64, submit func(func()) error) *TimingWheel {
	timingWheel := &TimingWheel{
		tickMs:             tickMs,
		status:             ETimingWheelStatusDead,
		slot:               [ETimerWheelDepth]int64{},
		bucket:             [ETimerWheelDepth][EBucketSize]*list.List{},
		overflow:           nil,
		deleteList:         nil,
		setCurrentSlotChan: make(chan int64, 1),
		addTaskChan:        make(chan *TaskData, EMaxConcurrentTaskNum),
		delTaskChan:        make(chan *TaskData, EMaxConcurrentTaskNum),
		exitChan:           make(chan struct{}),
		waitExitChan:       make(chan struct{}, 1),
		ticker:             nil,
		submit:             submit,
		currentSlot:        0,
		nextExpiry:         0,
	}
	for depth := 0; depth < ETimerWheelDepth; depth++ {
		for slot := 0; slot < EBucketSize; slot++ {
			timingWheel.bucket[depth][slot] = list.New()
		}
	}
	timingWheel.overflow = list.New()
	timingWheel.deleteList = list.New()
	return timingWheel
}

// GetTimingWheel 获取全局时间轮 建议全用全局
func GetTimingWheel() *TimingWheel {
	return gTimingWheel
}

// Init 初始化全局时间轮, submit 建议用协程池
func Init(tickMs int64, submit func(func()) error) {
	gTimingWheel = NewTimingWheel(tickMs, submit)
}
