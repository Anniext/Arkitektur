package websocket

import (
	"fmt"
	"github.com/Anniext/Arkitektur/system/log"
	"time"
)

var defaultWebsocket *WsServer

func InitDefaultWebsocket() error {
	cnf := GetDefaultWebsocketConfig()

	addr := fmt.Sprintf(":%d", cnf.Port)
	defaultWebsocket = NewWsServer(addr)

	// 暂时硬编码满足大部分情况， 如不能满足转到配置文件
	timeoutRead := time.Second * time.Duration(cnf.TimeoutRead)
	defaultWebsocket.WsSessionHub.SetTimeoutRead(timeoutRead)

	log.Info("server start in:", addr)
	go SafeGoRecoverWarpFunc(func() {
		defaultWebsocket.Start()
	})()

	return nil
}

func GetDefaultWebsocket() *WsServer {
	return defaultWebsocket
}
