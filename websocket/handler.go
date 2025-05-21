package websocket

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	server  *WsServer
	upgrade websocket.Upgrader
}

func (h *WsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if conn, err := h.upgrade.Upgrade(resp, req, nil); err != nil {
		return
	} else {
		wsSession := NewWsSession(&h.server.WsSessionHub, conn, req)
		h.server.sessions.Store(wsSession, true)
		atomic.AddInt32(&h.server.sessionNum, 1)
		log.Println("incoming connnetion ", wsSession.ClientIP())
	}
}
