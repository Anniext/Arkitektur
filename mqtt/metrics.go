package mqtt

// Metrics 客户端运行指标统计
type Metrics struct {
	Subscriptions    int64 // 当前订阅数
	MessagesSent     int64 // 已发送消息总数
	MessagesReceived int64 // 已接收消息总数
	HeartbeatsSent   int64 // 已发送心跳次数
	ConnectCount     int64 // 连接成功次数
	DisconnectCount  int64 // 断开连接次数
}
