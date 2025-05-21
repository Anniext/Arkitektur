package websocket

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var gServer *WsServer

// GetServer function    获取ws服务器配置
func GetServer() *WsServer {
	return gServer
}

// WsServer 服务器结构
type WsServer struct {
	httpServer    *http.Server // HTTP 服务器
	stopSignal    bool         // 停止信号
	exitFunctions []func()     // 退出回调
	WsSessionHub               // ws会话管理器
}

// NewWsServer function    新建ws服务器
func NewWsServer(addr string) *WsServer {
	server := &WsServer{}

	server.WsSessionHub.Init() // 初始化会话管理器
	server.httpServer = &http.Server{
		Addr: addr,
		Handler: &WsHandler{
			upgrade: websocket.Upgrader{
				HandshakeTimeout: 10 * time.Second,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			server: server,
		},
	}

	server.stopSignal = false
	gServer = server
	return server
}

// Start method    开启ws服务器
func (s *WsServer) Start() {
	var err error

	if s.certFile != "" && s.keyFile != "" {
		if err = s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil {
			log.Printf("https server close, err = %v", err)
		}

	} else {
		if err = s.httpServer.ListenAndServe(); err != nil {
			log.Printf("https server close, err = %v", err)
		}
	}

	sessions := make([]*WsSession, 1)

	s.sessions.Range(func(key *WsSession, value bool) bool {
		sessions = append(sessions, key)
		return true
	})

	// 清空已经有的连接
	for _, wsSession := range sessions {
		if wsSession != nil {
			wsSession.close(true)
		}
	}

	for _, fn := range s.exitFunctions {
		fn()
	}

	return
}

// SetTimeoutCloseRead method    设置关闭等待时间
func (s *WsServer) SetTimeoutCloseRead(timeout time.Duration) {
	s.timeoutCloseRead = timeout
}

// SetCert method    设置ssl证书
func (s *WsServer) SetCert(certFile, keyFile string) {
	s.certFile = certFile
	s.keyFile = keyFile
}

// SessionNum method   获取会话数
func (s *WsServer) SessionNum() int32 {
	return atomic.LoadInt32(&s.sessionNum)
}

// Stop method    停止服务
func (s *WsServer) Stop() {
	s.stopSignal = true
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Println("server stop error", err.Error())
	}
}

// AtClose method    添加关闭服务回调函数
func (s *WsServer) AtClose(fn func()) {
	s.exitFunctions = append(s.exitFunctions, fn)
}

// AtSessionClose method    添加会话注销函数
func (s *WsServer) AtSessionClose(fn func(*WsSession, int32)) {
	s.sessionExitFunctions = append(s.sessionExitFunctions, fn)
}
