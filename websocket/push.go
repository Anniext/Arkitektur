package websocket

import (
	"log"
	"sync/atomic"
	"time"
)

type WsSessionHub struct {
	sessions             *MapWsSessionBool                     // 建立连接的会话
	sessionNum           int32                                 // 建立连接的数量
	sessionMap           MapInt32WsSession                     // 并发map, key是msgNo, value是连接
	sessionExitFunctions []func(*WsSession, int32)             // 会话对应的回调
	protoFunctions       map[uint32]ProtoFunc                  // 协议对应的回调
	timeoutCloseRead     time.Duration                         // 关闭等待时间
	timeoutWrite         time.Duration                         // 写包超时时间
	timeoutRead          time.Duration                         // 读包超时时间
	message              IMessage                              // 包结构
	UnregisteredCallback ProtoFunc                             // 未注册的回调函数
	middleware           []func(protoFunc ProtoFunc) ProtoFunc // 中间件
	certFile             string                                // 安全证书
	keyFile              string                                // 安全密钥
	ForwardedByClientIP  bool                                  // 白名单
}

// Init 初始化
func (hub *WsSessionHub) Init() {
	hub.protoFunctions = make(map[uint32]ProtoFunc)
	hub.sessions = &MapWsSessionBool{}
	hub.timeoutCloseRead = 0
	hub.timeoutWrite = 5 * time.Minute
	hub.timeoutRead = 5 * time.Minute
	hub.message = &Message{}
	hub.ForwardedByClientIP = true
}

// Use method    使用中间件
func (hub *WsSessionHub) Use(middleware func(protoFunc ProtoFunc) ProtoFunc) {
	hub.middleware = append(hub.middleware, middleware)
}

// Register method    注册
func (hub *WsSessionHub) Register(msgNo uint32, fn ProtoFunc) {
	if fn == nil {
		return
	}

	_, ok := hub.protoFunctions[msgNo]
	if ok {
		return
	} else {
		for idx := range hub.middleware {
			fn = hub.middleware[idx](fn)
			hub.protoFunctions[msgNo] = fn
		}
	}
}

// Exit method    注册
func (hub *WsSessionHub) Exit() {
	sessions := make([]*WsSession, 0)
	hub.sessions.Range(func(key *WsSession, value bool) bool {
		sessions = append(sessions, key)
		return true
	})

	for _, ws := range sessions {
		if ws != nil {
			ws.Close()
		}
	}

	bTime := time.Now().Unix()
	waitSec := int64(hub.timeoutRead.Seconds())

	for {
		sessions := make([]*WsSession, 0)
		hub.sessions.Range(func(key *WsSession, value bool) bool {
			sessions = append(sessions, key)
			return true
		})

		var exitNum int

		for _, wsSession := range sessions {
			if atomic.LoadInt32(&wsSession.state) >= sessionExit {
				exitNum += 1
			}
		}

		if exitNum >= len(sessions) {
			log.Println("session hub exit ", len(sessions), exitNum)
			break
		}

		if time.Now().Unix() > bTime+waitSec {
			log.Println("wait session exit timeout", len(sessions), exitNum)
			break
		}

		time.Sleep(time.Second)
	}
}

// Get method    获取会话
func (hub *WsSessionHub) Get(uid int32) (*WsSession, bool) {
	v, has := hub.sessionMap.Load(uid)

	if has {
		session := v
		return session, true
	}

	return nil, false
}

// GetSessionAll method    获取所有的会话
func (hub *WsSessionHub) GetSessionAll() (sessionList []*WsSession) {
	hub.sessionMap.Range(func(key int32, value *WsSession) bool {
		sessionList = append(sessionList, value)
		return true
	})

	return
}

// AddSession method    增加会话
func (hub *WsSessionHub) AddSession(uid int32, s *WsSession) (actual *WsSession, loaded bool) {
	actual, loaded = hub.sessionMap.LoadOrStore(uid, s)

	return
}

// AtSessionClose method    增加会话关闭回调函数
func (hub *WsSessionHub) AtSessionClose(fn func(hub *WsSession, i int32)) {
	hub.sessionExitFunctions = append(hub.sessionExitFunctions, SafeGoRecoverWarpFuncUid(fn))
}

// RemoveSession method    移除会话
func (hub *WsSessionHub) RemoveSession(uid int32) {
	hub.sessionMap.Delete(uid)
}

// SessionNum method    获取会话数量
func (hub *WsSessionHub) SessionNum() int32 {
	return atomic.LoadInt32(&hub.sessionNum)
}

// SetTimeoutCloseRead method    设置超过关闭等待时间
func (hub *WsSessionHub) SetTimeoutCloseRead(timeout time.Duration) {
	hub.timeoutCloseRead = timeout
}

// SetTimeoutWrite method    设置写入超时时间
func (hub *WsSessionHub) SetTimeoutWrite(timeout time.Duration) {
	hub.timeoutWrite = timeout
}

// SetTimeoutRead method    设置读取超时时间
func (hub *WsSessionHub) SetTimeoutRead(timeout time.Duration) {
	hub.timeoutRead = timeout
}

// SetMessage method    设置会话
func (hub *WsSessionHub) SetMessage(message IMessage) {
	hub.message = message
}

// GetSession method    获取会话
func (hub *WsSessionHub) GetSession(uid int32) *WsSession {
	session, ok := hub.sessionMap.Load(uid)
	if ok {
		return session
	}
	return nil
}

// PushMsg method    向会话推送
func (hub *WsSessionHub) PushMsg(uid int32, message IMessage) bool {
	session, ok := hub.sessionMap.Load(uid)
	if !ok {
		return false
	}
	return session.PushMsg(message)
}

// PushWork method    向工作队列推送
func (hub *WsSessionHub) PushWork(uid int32, fn func()) bool {
	s, ok := hub.sessionMap.Load(uid)
	if !ok {
		return false
	}
	return s.PushWork(fn)
}

// PushToAll method    向所有ws推送消息
func (hub *WsSessionHub) PushToAll(message IMessage) {
	var all []*WsSession
	hub.sessionMap.Range(func(key int32, value *WsSession) bool {
		all = append(all, value)
		return true
	})

	for _, s := range all {
		s.PushMsg(message)
	}
}

func (hub *WsSessionHub) DoSomething(fn func(uid int32, session *WsSession) bool) {
	hub.sessionMap.Range(func(key int32, value *WsSession) bool {
		return fn(key, value)
	})
}
