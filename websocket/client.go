package websocket

import (
	"errors"
	"log"
	"net/url"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const EMaxReconnect = 3
const EReconnectDelay = 5 * time.Millisecond
const EMaxRecheck = 3
const EDefaultClientNo = 1

func getSidEasy(session *WsSession, message IMessage) uint64 {
	if message != nil {
		//return uint64()
		return uint64(EDefaultClientNo)<<32 + uint64(message.GetMsgNo())
	}
	return 0
}

type WsClient struct {
	WsSessionHub
	currentSessionNo     int32                                             // 当前会话编号
	clientNo2SessionNo   []int32                                           // 客户端变化对应会话编号
	clientAddresses      []string                                          // 客户端链接地址，更换会导致并发问题
	currentClientNo      int32                                             // 当前有效的客户端编号
	maxClientNo          int32                                             // 最大客户端编号， 创建后无法修改
	msgReceive           MapUint64ChanBytes                                // 消息接受通道
	timeoutWait          time.Duration                                     // 协议回包超时时间，超时后会话会关闭
	getSidCallback       func(session *WsSession, message IMessage) uint64 // 获取sid
	clientIdPool         []int32                                           // 回收的clientId 池子
	clientIdPoolMutex    sync.Mutex                                        // clientId 池子 变化锁
	connectDoneFunctions []func(int32)                                     // 链接完成后回调
	clientReconnectState []int32                                           // 客户端重连状态 避免并发
}

func NewWsClient(maxClientNo int32, getSid func(session *WsSession, message IMessage) uint64) *WsClient {
	return NewWsClientWittCoder(maxClientNo, &Message{}, getSid)
}

func NewWsClientWittCoder(maxClientNo int32, message IMessage, getSid func(session *WsSession, message IMessage) uint64) *WsClient {
	client := new(WsClient)
	client.WsSessionHub.Init()
	client.WsSessionHub.UnregisteredCallback = client.UnregisteredCallback
	client.message = message
	client.timeoutWait = 30 * time.Second
	client.clientNo2SessionNo = make([]int32, maxClientNo+1)
	client.clientAddresses = make([]string, maxClientNo+1)
	client.maxClientNo = maxClientNo
	client.getSidCallback = getSid
	client.clientReconnectState = make([]int32, maxClientNo+1)
	client.ForwardedByClientIP = false
	return client
}

func NewWsClientEasy(addr string) *WsClient {
	return NewWsClientEasyWittCoder(addr, &Message{})
}

func NewWsClientEasyWittCoder(addr string, coder IMessage) *WsClient {
	client := NewWsClientWittCoder(EDefaultClientNo, coder, getSidEasy)
	clientNo := client.Connect(addr)
	if clientNo < 0 {
		client.WsSessionHub.Exit()
		return nil
	}
	return client

}

func (c *WsClient) Stop() {
	c.WsSessionHub.Exit()
	return
}

// get SessionNo
func (c *WsClient) GetSessionNo(clientNo int32) (sessionNo int32) {
	if clientNo > c.maxClientNo || clientNo < 0 {
		log.Println("invalid clientNo ", clientNo)
		return -1
	}
	sessionNo = atomic.LoadInt32(&c.clientNo2SessionNo[clientNo])
	return
}

func (c *WsClient) GetMaxClientNo() int32 {
	return c.maxClientNo
}

func (c *WsClient) SetTimeoutWait(timeout time.Duration) {
	c.timeoutWait = timeout
}

func (c *WsClient) UnregisteredCallback(session *WsSession, message IMessage) []byte {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			log.Println("stack:", string(debug.Stack()))
		} else {
		}
	}()
	sid := c.getSidCallback(session, message)
	if ch, ok := c.msgReceive.Load(sid); ok {
		select {
		case ch <- message.GetBody():
		default:
			log.Println("UnregisteredCallback error", session.GetUid(), sid, message.GetMsgNo())
		}
	}
	return nil
}

func (c *WsClient) AtConnectDone(fn func(int322 int32)) {
	c.connectDoneFunctions = append(c.connectDoneFunctions, fn)
}

func (c *WsClient) Connect(addr string) (clientNo int32) {

	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	log.Println(u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		log.Println("dail fail:", addr)
		return -1
	}
	session := NewWsSession(&c.WsSessionHub, conn, nil)
	if session == nil {
		log.Println("NewTcpSession fail")
		return -1
	}
	// 检查回收池
	c.clientIdPoolMutex.Lock()
	length := len(c.clientIdPool)
	if length > 0 {
		clientNo = c.clientIdPool[length-1]
		c.clientIdPool = c.clientIdPool[0 : length-1]
	}
	c.clientIdPoolMutex.Unlock()

	// 暂无回收client no
	if clientNo <= 0 {
		var currentClientNo int32
		for i := 0; i < EMaxRecheck; i++ {
			currentClientNo = atomic.LoadInt32(&c.currentClientNo)
			if currentClientNo >= c.maxClientNo {
				log.Println("too many clients ", currentClientNo)
				session.close(false)
				return -1
			}
			if atomic.CompareAndSwapInt32(&c.currentClientNo, currentClientNo, currentClientNo+1) {
				clientNo = currentClientNo + 1
				break
			}
		}
	}
	// clientNo 有效
	if clientNo > 0 {
		sessionNo := atomic.AddInt32(&c.currentSessionNo, 1)
		success := session.SetUid(sessionNo)
		if success {
			atomic.StoreInt32(&c.clientNo2SessionNo[clientNo], sessionNo)
			c.clientAddresses[clientNo] = addr
			for _, fn := range c.connectDoneFunctions {
				fn(clientNo)
			}
			return
		}
		c.clientIdPoolMutex.Lock()
		c.clientIdPool = append(c.clientIdPool, clientNo)
		c.clientIdPoolMutex.Unlock()
	}
	session.close(false)
	return -1
}

func (c *WsClient) Reconnect(clientNo int32) (sessionNo int32) {
	if clientNo > c.maxClientNo || clientNo < 0 {
		log.Println("invalid clientNo ", clientNo)
		return -1
	}
	if atomic.CompareAndSwapInt32(&c.clientReconnectState[clientNo], 0, 1) {
		defer func() {
			atomic.CompareAndSwapInt32(&c.clientReconnectState[clientNo], 1, 0)

		}()
		// 其他人已重新连接
		sessionNo = atomic.LoadInt32(&c.clientNo2SessionNo[clientNo])
		if sessionNo > 0 {
			session := c.WsSessionHub.GetSession(sessionNo)
			if session != nil && !session.Dead() {
				return sessionNo
			}
		}
		addr := c.clientAddresses[clientNo]
		c.CloseSessionByClientNo(clientNo)
		atomic.StoreInt32(&c.clientNo2SessionNo[clientNo], -1)
		conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
		if err != nil {
			log.Println("dial fail:", addr)
			return -1
		}
		session := NewWsSession(&c.WsSessionHub, conn, nil)
		if session == nil {
			log.Println("NewWsSession fail")
			return -1
		}
		sessionNo = atomic.AddInt32(&c.currentSessionNo, 1)
		success := session.SetUid(sessionNo)
		if !success {
			session.close(false)
			return -1
		}
		atomic.StoreInt32(&c.clientNo2SessionNo[clientNo], sessionNo)
		for _, fn := range c.connectDoneFunctions {
			fn(clientNo)
		}
	}

	return -1

}

// CloseSessionByClientNo 断开客户端session
func (c *WsClient) CloseSessionByClientNo(clientNo int32) {
	sessionNo := c.GetSessionNo(clientNo)
	session := c.WsSessionHub.GetSession(sessionNo)
	if session != nil {
		session.close(false)
	}
}

// Disconnect 关闭客户端连接回收资源
func (c *WsClient) Disconnect(clientNo int32) {
	c.CloseSessionByClientNo(clientNo)
	// 回收clientNo
	if clientNo > 0 {
		c.clientIdPoolMutex.Lock()
		c.clientIdPool = append(c.clientIdPool, clientNo)
		c.clientIdPoolMutex.Unlock()
	}
}

// 检查连接状态
func (c *WsClient) checkClientConnection(clientNo int32) (sessionNo int32, err error) {
	needReconnect := false
	if clientNo < 0 {
		return 0, errors.New("invalid clientNo")
	}

	sessionNo = c.GetSessionNo(clientNo)
	if sessionNo <= 0 {
		needReconnect = true
	}
	session := c.WsSessionHub.GetSession(sessionNo)

	if session != nil || session.Dead() {
		needReconnect = true
	}

	if needReconnect {
		for i := 0; i < EMaxReconnect; i++ {
			sessionNo = c.Reconnect(clientNo)
			if sessionNo > 0 {
				break
			}
			time.Sleep(EReconnectDelay)
		}
	}
	return sessionNo, nil
}

func (c *WsClient) Request(clientNo int32, sid uint64, message IMessage) (resp []byte, err error) {

	sessionNo, err := c.checkClientConnection(clientNo)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte, 1)

	c.msgReceive.Store(sid, ch)
	c.WsSessionHub.PushMsg(sessionNo, message)

	select {
	case tmp, ok := <-ch:
		if ok {
			resp = tmp
		}
	case <-time.After(c.timeoutWait):
		log.Println("timeout", clientNo, sid, message.GetMsgNo())
		c.CloseSessionByClientNo(clientNo)
		err = errors.New("timeout")
	}
	c.msgReceive.Delete(sid)

	return
}

func (c *WsClient) RequestNoWait(clientNo int32, sid uint64, message IMessage) error {

	sessionNo, err := c.checkClientConnection(clientNo)
	if err != nil {
		return err
	}
	c.WsSessionHub.PushMsg(sessionNo, message)
	return nil
}

func (c *WsClient) RequestEasy(message IMessage) ([]byte, error) {
	return c.Request(
		EDefaultClientNo,
		uint64(EDefaultClientNo)<<32+uint64(message.GetMsgNo()),
		message,
	)
}
