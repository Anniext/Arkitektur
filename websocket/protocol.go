package websocket

type ProtoFunc func(*WsSession, IMessage) []byte
