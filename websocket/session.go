package websocket

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	sessionRunning = iota
	sessionStop
	sessionExit
)

type WsSession struct {
	uid          int32
	state        int32
	conn         *websocket.Conn // ws连接句柄
	writeQueue   *MessageQueue   // 消息队列
	msg          *Message        //消息
	workQueue    *WorkQueue      // 工作队列
	hub          *WsSessionHub   //会话管理
	Request      *http.Request   // http请求句柄
	wg           sync.WaitGroup
	timeoutRead  time.Duration // 读超时
	timeoutWrite time.Duration // 写超时
	Context      sync.Map
}

// NewWsSession function    新建ws会话管理
func NewWsSession(hub *WsSessionHub, conn *websocket.Conn, req *http.Request) *WsSession {
	ws := &WsSession{}

	ws.conn = conn
	ws.Request = req
	ws.writeQueue = NewMessageQueue()
	ws.workQueue = NewWorkQueue()
	ws.wg.Add(2)
	ws.timeoutRead = 5 * time.Minute
	ws.timeoutWrite = 5 * time.Minute
	ws.hub = hub

	hub.sessions.Store(ws, true)
	atomic.AddInt32(&hub.sessionNum, 1)

	coder := hub.message
	remoteAddr := ws.ClientIP()
	log.Println("new session: ", remoteAddr)

	// 循环从写队列中取出然后发送给ws链接
	go SafeGoRecoverWarpFunc(func() {
		var writeList [][]byte
		var err error

	LabelWriteThread:
		for {
			writeList := writeList[0:0]
			exit := ws.writeQueue.Pick(&writeList)

			for _, msg := range writeList {
				err = ws.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					break LabelWriteThread
				}
			}

			if exit {
				break
			}
		}

	})()

	go SafeGoRecoverWarpFunc(func() {
		// 处理正在进行中的工作
		for {
			if atomic.LoadInt32(&ws.state) >= sessionStop {
				log.Println("recv thread session close", remoteAddr)
				break
			}

			readBeginTime := time.Now()
			conn.SetReadDeadline(time.Now().Add(ws.timeoutRead))

			msg, err := coder.Decode(conn)
			if err != nil {
				var netErr *net.OpError
				if errors.As(err, &netErr) && netErr.Timeout() {
					if time.Now().Sub(readBeginTime) >= ws.timeoutRead {
						log.Println("recv thread timeout", remoteAddr, err)
						// todo 超时错误
						continue
					}
				} else {
					// TODO 其他错误解决
					break
				}
			} else if msg == nil {
				log.Println("recv thread msg is nil", err)
			} else {
				// 收到消息
				fn, ok := ws.hub.protoFunctions[msg.GetMsgNo()]
				if ok {
					SafeGoRecoverWarpFunc(func() {
						//log.Println("handler begin: ", remoteAddr)

						outMsg := &Message{
							MsgNo: msg.GetMsgNo(),
							Body:  fn(ws, msg),
						}
						//log.Println("handler end: ", remoteAddr)
						outMsg.SetLength()

						if outMsg.Body != nil {
							msgStr, err := coder.Encode(outMsg)
							if err != nil {
								log.Println("proto handler", remoteAddr, err)
							} else {
								// 将协议发给对应的回调
								ws.sendMsg(msgStr)
							}
						}

						if ws.msg != nil {
							msgStr, err := coder.Encode(ws.msg)
							if err != nil {
								log.Println("post msg: ", remoteAddr, err)
							} else {
								ws.sendMsg(msgStr)
							}
							ws.msg = nil
						}
					})()

				} else {
					if ws.hub.UnregisteredCallback == nil {
						log.Println("unkonw msg no", remoteAddr, msg.GetMsgNo())
						break
					} else {
						// 没有回调的时候走未注册的回调
						ws.hub.UnregisteredCallback(ws, msg)
					}
				}
			}

			// 从工作队列里面拿任务然后工作
			funcList := ws.workQueue.Dump()
			for _, work := range funcList {
				work()
			}
		}

		// 处理没有完成的工作
		funcList := ws.workQueue.Dump()
		for _, work := range funcList {
			work()
		}

		// 处理会话退出的任务
		for _, fn := range ws.hub.sessionExitFunctions {
			fn(ws, ws.GetUid())
		}
	})()

	return ws
}

// GetTimeoutRead method    获取读取超时时间
func (s *WsSession) GetTimeoutRead() time.Duration {
	return s.timeoutRead
}

// SetTimeoutRead method    设置读取超时时间
func (s *WsSession) SetTimeoutRead(timeout time.Duration) {
	s.timeoutRead = timeout
}

// GetTimeoutWrite   获取写入超时时间
func (s *WsSession) GetTimeoutWrite() time.Duration {
	return s.timeoutWrite
}

// SetTimeoutWrite   设置写入读取时间
func (s *WsSession) SetTimeoutWrite(writeout time.Duration) {
	s.timeoutWrite = writeout
}

// SetReadDeadline  设置读取等待时间
func (s *WsSession) SetReadDeadline(t time.Time) error {
	return s.conn.SetReadDeadline(t)
}

// close method    关闭会话
func (s *WsSession) close(wait bool) bool {
	if atomic.CompareAndSwapInt32(&s.state, sessionRunning, sessionStop) {
		s.workQueue.Add(nil)
		s.hub.sessions.Delete(s)
		if s.uid != 0 {
			s.hub.RemoveSession(s.uid)
		}

		s.conn.SetReadDeadline(time.Now().Add(s.hub.timeoutCloseRead))
		if wait {
		}
	}

	return true
}

// Close    关闭会话
func (s *WsSession) Close() bool {
	return s.close(true)
}

// Dead method    会话是否结束
func (s *WsSession) Dead() bool {
	return atomic.LoadInt32(&s.state) >= sessionStop
}

// RemoteAddr method    获取远程ip
func (s *WsSession) RemoteAddr() string {
	if s != nil && s.conn != nil {
		return s.conn.RemoteAddr().String()
	}

	return ""
}

// ClientIP method    获取ip地址
func (s *WsSession) ClientIP() string {
	if s.Request == nil || !s.hub.ForwardedByClientIP {
		return s.RemoteAddr()
	}

	if s.hub.ForwardedByClientIP {
		clientIP := s.Request.Header.Get("X-Forwarded-For")
		log.Println(clientIP)

		if len(clientIP) == 0 {
			clientIP = strings.TrimSpace(s.Request.Header.Get("X-Real-Ip"))
		}

		port := s.Request.Header.Get("X-Real-Port")
		if clientIP != "" {
			return clientIP + ":" + port
		}

	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(s.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// sendMsg method    将消息发送到工作队列里面
func (s *WsSession) sendMsg(buf []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", s)
		}
	}()

	if atomic.LoadInt32(&s.state) == sessionRunning {
		s.writeQueue.Add(buf)
		return true
	} else {
		return false
	}
}

// GetUid method    获取uid
func (s *WsSession) GetUid() int32 {
	return s.uid
}

// SetUid method    设置uid
func (s *WsSession) SetUid(uid int32) bool {
	_, loaded := s.hub.AddSession(uid, s)
	if !loaded {
		s.uid = uid
		if s.Dead() {
			s.uid = 0
			s.hub.RemoveSession(uid)
			return false
		}
		return true
	} else {
		return false
	}
}

// SetUidSafe method    安全设置uid
func (s *WsSession) SetUidSafe(uid int32) bool {
	_, loaded := s.hub.AddSession(uid, s)
	if !loaded {
		atomic.StoreInt32(&s.uid, 0)
		if s.Dead() {
			atomic.StoreInt32(&s.uid, 0)
			s.hub.RemoveSession(uid)
			return false
		}
		return true
	} else {
		return false
	}
}

// GetUidSafe method    获取uid
func (s *WsSession) GetUidSafe() int32 {
	return atomic.LoadInt32(&s.uid)
}

// PushMsg method    将消息编码字节后发送
func (s *WsSession) PushMsg(message IMessage) bool {
	msgStr, err := s.hub.message.Encode(message)
	if err != nil {
		return false
	}
	return s.sendMsg(msgStr)
}

// PushWork method    将任务函数提交到工作队列
func (s *WsSession) PushWork(fn func()) bool {
	if atomic.LoadInt32(&s.state) == sessionRunning {
		s.workQueue.Add(SafeGoRecoverWarpFunc(fn))
		return true
	}
	return false
}

// Post method    提交消息到会话管理器
func (s *WsSession) Post(msgNo uint32, content []byte) {
	s.msg = &Message{
		MsgNo: msgNo,
		Body:  content,
	}
	s.msg.SetLength()
}

// waitAndClose method    等待并关闭连接
func (s *WsSession) waitAndClose() {
	SafeGoRecoverWarpFunc(func() {
		s.wg.Wait()
		atomic.AddInt32(&s.hub.sessionNum, -1)
		s.hub.sessions.Delete(s)
		if s.uid != 0 {
			s.hub.RemoveSession(s.uid)
		}

		s.workQueue.Reset()
		s.writeQueue.Reset()
		s.conn.Close()
	})
}
