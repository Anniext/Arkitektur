package mqtt

import (
	"admin/common/log"
	"admin/common/recover"
	"context"
	"encoding/json"
	"fmt"
	mqttx "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"sync/atomic"
	"time"
)

// MQTTClient MQTT客户端结构体
type MQTTClient struct {
	client        mqttx.Client
	config        ClientConfig
	mu            sync.RWMutex
	connected     bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	metrics       Metrics
	subscriptions sync.Map // 保存主题与回调的映射
	pool          chan struct{}
	stopOnce      sync.Once
	Chan          chan struct{}
}

// NewClient 创建新的MQTT客户端实例
// 参数：
//
//	cfg - 客户端配置
//	handler - 消息接收处理回调函数
//
// 返回：
//
//	*MQTTClient 初始化的客户端实例
func NewClient(cfg ClientConfig) *MQTTClient {
	if cfg.ClientID == "" {
		cfg.ClientID = generateClientID()
	}
	if cfg.KeepAlive == 0 {
		cfg.KeepAlive = 30 * time.Second
	}
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}
	if cfg.MaxReconnectInterval == 0 {
		cfg.MaxReconnectInterval = 5 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &MQTTClient{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
		pool:   make(chan struct{}, 20),
		Chan:   make(chan struct{}, 1),
	}
}

// Connect 连接到MQTT代理服务器
// 返回：
//
//	error 连接错误信息
func (c *MQTTClient) Connect() error {
	opts := mqttx.NewClientOptions()

	// 设置基础参数
	opts.AddBroker(c.config.BrokerURL).
		SetClientID(c.config.ClientID).
		SetKeepAlive(c.config.KeepAlive).
		SetConnectTimeout(c.config.ConnectTimeout).
		SetAutoReconnect(c.config.AutoReconnect).
		SetMaxReconnectInterval(c.config.MaxReconnectInterval).
		SetOnConnectHandler(c.onConnect).             // 连接成功回调
		SetConnectionLostHandler(c.onConnectionLost). // 连接丢失回调
		SetReconnectingHandler(c.onReconnecting)      // 重连中回调

	// 配置TLS
	if c.config.TLSConfig != nil {
		opts.SetTLSConfig(c.config.TLSConfig)
	}

	// 设置遗嘱消息
	if c.config.WillTopic != "" {
		opts.SetWill(c.config.WillTopic, string(c.config.WillPayload),
			c.config.WillQoS, c.config.WillRetain)
	}

	// 创建客户端实例
	c.client = mqttx.NewClient(opts)

	// 执行连接
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("连接失败: %w", token.Error())
	}

	// 启动后台协程
	c.wg.Add(1)

	go recoverutils.RecoverWrapFunc(func() {
		c.heartbeatLoop() // 心跳协程
	})()

	return nil
}

// Disconnect 断开连接并释放资源
func (c *MQTTClient) Disconnect() {
	c.stopOnce.Do(func() {
		// 通知所有协程停止
		c.cancel()

		// 断开MQTT连接
		c.client.Disconnect(250)

		// 等待所有协程退出
		c.wg.Wait()

		// 关闭协程池
		close(c.pool)

		log.Infof("客户端已完全关闭")
	})
}

// Publish 发布消息到指定主题
// 参数：
//
//	topic - 消息主题
//	payload - 消息内容
//
// 返回：
//
//	error 发布错误信息
func (c *MQTTClient) Publish(topic string, payload []byte) error {
	if !c.IsConnected() {
		return fmt.Errorf("客户端未连接")
	}

	token := c.client.Publish(topic, c.config.QoS, c.config.Retain, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	atomic.AddInt64(&c.metrics.MessagesSent, 1)
	return nil
}

// Subscribe 订阅指定主题
// 参数：
//
//	topic - 要订阅的主题
//
// 返回：
//
//	error 订阅错误信息
func (c *MQTTClient) Subscribe(topic string, callback func(string, []byte)) error {
	if !c.IsConnected() {
		return fmt.Errorf("客户端未连接")
	}
	// 保存订阅关系
	c.subscriptions.Store(topic, callback)
	atomic.AddInt64(&c.metrics.Subscriptions, 1)

	// 注册消息路由
	c.client.AddRoute(topic, func(_ mqttx.Client, msg mqttx.Message) {
		atomic.AddInt64(&c.metrics.MessagesReceived, 1)
		callback(msg.Topic(), msg.Payload())
	})

	// 发送订阅请求
	if token := c.client.Subscribe(topic, c.config.QoS, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Infof("成功订阅主题: %s", topic)
	return nil
}

// Unsubscribe 取消订阅主题
// 参数：
//
//	topics - 要取消订阅的主题列表
//
// 返回：
//
//	error 取消订阅错误信息
func (c *MQTTClient) Unsubscribe(topics ...string) error {
	if !c.IsConnected() {
		return fmt.Errorf("客户端未连接")
	}

	for _, topic := range topics {
		if _, loaded := c.subscriptions.LoadAndDelete(topic); loaded {
			atomic.AddInt64(&c.metrics.Subscriptions, -1)
		}
	}

	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// IsConnected 返回当前连接状态
func (c *MQTTClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// GetMetrics 获取当前运行指标
func (c *MQTTClient) GetMetrics() Metrics {
	return Metrics{
		Subscriptions:    atomic.LoadInt64(&c.metrics.Subscriptions),
		MessagesSent:     atomic.LoadInt64(&c.metrics.MessagesSent),
		MessagesReceived: atomic.LoadInt64(&c.metrics.MessagesReceived),
		HeartbeatsSent:   atomic.LoadInt64(&c.metrics.HeartbeatsSent),
		ConnectCount:     atomic.LoadInt64(&c.metrics.ConnectCount),
		DisconnectCount:  atomic.LoadInt64(&c.metrics.DisconnectCount),
	}
}

// 心跳循环
func (c *MQTTClient) heartbeatLoop() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.config.KeepAlive)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.sendHeartbeat()
		case <-c.ctx.Done():
			return
		}
	}
}

// 发送心跳消息
func (c *MQTTClient) sendHeartbeat() {
	alive := struct {
		Status   string `json:"status"`
		ClientID string `json:"client_id"`
		Ts       int64  `json:"ts"`
	}{Status: "alive", ClientID: c.config.ClientID, Ts: time.Now().Unix()}
	payload, _ := json.Marshal(alive)

	if err := c.Publish("device/status", payload); err != nil {
		log.Errorf("心跳发送失败: %v", err)
		return
	}
	atomic.AddInt64(&c.metrics.HeartbeatsSent, 1)
}

// 连接成功回调
func (c *MQTTClient) onConnect(_ mqttx.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = true

	// 重新订阅所有主题
	c.subscriptions.Range(func(topic, callback interface{}) bool {
		t := topic.(string)
		c.client.Subscribe(t, c.config.QoS, nil)
		return true
	})

	c.Chan <- struct{}{} // 通知所有协程订阅
	log.Infof("连接成功!!!")
}

// 连接丢失回调
func (c *MQTTClient) onConnectionLost(_ mqttx.Client, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
	atomic.AddInt64(&c.metrics.DisconnectCount, 1)
	log.Errorf("连接丢失: %v", err)
}

// 重连中回调
func (c *MQTTClient) onReconnecting(_ mqttx.Client, _ *mqttx.ClientOptions) {
	log.Infof("尝试重新连接...")
}

// 指标收集协程
func (c *MQTTClient) metricsCollector() {
	defer c.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := c.GetMetrics()
			log.Debugf("运行指标 - 发送: %d, 接收: %d, 心跳: %d",
				metrics.MessagesSent, metrics.MessagesReceived, metrics.HeartbeatsSent)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *MQTTClient) GetMqttClientID() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.ClientID
}
