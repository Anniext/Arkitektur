package timer

import "container/list"

// TaskData 任务实例
type TaskData struct {
	Element         *list.Element // element
	List            *list.List    // list
	ExpiresSlot     int64         // 过期槽数
	ExpiresLoopSlot int64         // 循环槽数
	Callback        func()        // 回调
}
