package websocket

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

const MaxPacketLen uint32 = 1024 * 1024 * 1000

type IMessage interface {
	GetMsgNo() uint32
	GetLength() uint32
	SetLength()
	SetBody([]byte)
	GetBody() []byte
	Decode(*websocket.Conn) (IMessage, error)
	Encode(IMessage) ([]byte, error)
}

type Message struct {
	MsgNo  uint32
	Length uint32
	Body   []byte
}

func (p *Message) GetMsgNo() uint32 {
	if p == nil {
		return 0
	}

	return p.MsgNo
}

func (p *Message) GetLength() uint32 {
	if p == nil {
		return 0
	}

	return p.Length
}

func (p *Message) SetLength() {
	if p == nil {
		return
	}

	p.Length = uint32(8) + uint32(len(p.Body))
}

func (p *Message) SetBody(bytes []byte) {
	if p == nil {
		return
	}

	p.Body = bytes
}

func (p *Message) GetBody() []byte {
	if p == nil {
		return nil
	}

	return p.Body
}

func (p *Message) Decode(conn *websocket.Conn) (IMessage, error) {
	_, dataReader, err := conn.ReadMessage()

	if err != nil {
		// todo 调试
		return nil, err
	}

	pm := &Message{}

	// 前四位是消息号
	pm.MsgNo = uint32(dataReader[0]&0xff) + uint32(dataReader[1]&0xff)<<8 + uint32(dataReader[2]&0xff)<<16 + uint32(dataReader[3]&0xff)<<24

	// 后两位是长度
	pm.Length = uint32(dataReader[4]&0xff) + uint32(dataReader[5]&0xff)<<8 + uint32(dataReader[6]&0xff)<<16 + uint32(dataReader[7]&0xff)<<24
	if pm.Length > MaxPacketLen {

		return nil, errors.New(fmt.Sprintf("large packet %d", pm.Length))
	}

	// 后面的是消息
	pm.Body = dataReader[8:]

	return pm, nil
}

func (p *Message) Encode(msg IMessage) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}

	var body []byte
	body = msg.GetBody()

	var l = len(body)
	var buf [4]byte
	var msgOut []byte

	binary.LittleEndian.PutUint32(buf[:], msg.GetMsgNo())
	msgOut = append(msgOut, buf[:]...)

	binary.LittleEndian.PutUint32(buf[:], uint32(l+8))
	msgOut = append(msgOut, buf[:]...)

	msgOut = append(msgOut, body...)
	return msgOut, nil
}
