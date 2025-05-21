//go:generate rm -rf ./001_*
//go:generate syncmap -name=MapWsSessionBool -pkg=websocket -o=./001_mapwssessionbool.go map[*WsSession]bool
//go:generate syncmap -name=MapInt32WsSession -pkg=websocket -o=./001_mapint32wssession.go map[int32]*WsSession
//go:generate syncmap -name=MapUint64ChanBytes -pkg=websocket -o=./001_mapuint64chanbytes.go map[uint64]chan[]byte
package websocket
