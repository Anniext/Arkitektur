package mqtt

type MqttConfig struct {
	BrokerURL            string
	ClientID             string
	KeepAlive            int
	QoS                  byte
	Retain               bool
	AutoReconnect        bool
	ConnectTimeout       int
	MaxReconnectInterval int
	ServerID             string
}
type Option func(*MqttConfig)

func WithBrokerURLOption(brokerURL string) Option {
	return func(c *MqttConfig) {
		c.BrokerURL = brokerURL
	}
}
func WithClientIDOption(clientID string) Option {
	return func(c *MqttConfig) {
		c.ClientID = clientID
	}
}
func WithKeepAliveOption(keepAlive int) Option {
	return func(c *MqttConfig) {
		c.KeepAlive = keepAlive
	}
}
func WithQoSOption(qoS byte) Option {
	return func(c *MqttConfig) {
		c.QoS = qoS
	}
}
func WithRetainOption(retain bool) Option {
	return func(c *MqttConfig) {
		c.Retain = retain
	}
}
func WithAutoReconnectOption(autoReconnect bool) Option {
	return func(c *MqttConfig) {
		c.AutoReconnect = autoReconnect
	}
}
func WithConnectTimeoutOption(connectTimeout int) Option {
	return func(c *MqttConfig) {
		c.ConnectTimeout = connectTimeout
	}
}
func WithMaxReconnectIntervalOption(maxReconnectInterval int) Option {
	return func(c *MqttConfig) {
		c.MaxReconnectInterval = maxReconnectInterval
	}
}

func WithServerIDOption(serverID string) Option {
	return func(c *MqttConfig) {
		c.ServerID = serverID
	}
}

func NewMqttOption(options ...Option) {
	defaultMqttConfig = &MqttConfig{}
	for _, option := range options {
		option(defaultMqttConfig)
	}
}

var defaultMqttConfig *MqttConfig

func GetDefaultMqttConfig() *MqttConfig {
	return defaultMqttConfig
}
