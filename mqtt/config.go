package mqtt

import (
	"crypto/tls"
	"time"
)

type ClientConfig struct {
	BrokerURL            string        // mqtt服务器
	ClientID             string        // 客户端id
	KeepAlive            time.Duration // 心跳包
	QoS                  byte          // 消息质量
	Retain               bool          // 保留消息标志
	AutoReconnect        bool          // 是否自重连
	ConnectTimeout       time.Duration // 连接超时时间
	MaxReconnectInterval time.Duration // 最大重连间隔
	WillTopic            string        // 遗嘱消息主题
	WillPayload          []byte        // 遗嘱消息内容
	WillQoS              byte          // 遗嘱消息质量
	WillRetain           bool          // 遗嘱保留标志
	TLSConfig            *tls.Config   // TLS配置
}
