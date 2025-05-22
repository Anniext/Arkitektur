package websocket

type WebsocketConfig struct {
	Port        int
	TimeoutRead int
}
type Option func(*WebsocketConfig)

func WithPortOption(port int) Option {
	return func(c *WebsocketConfig) {
		c.Port = port
	}
}

func WithTimeoutReadOption(timeoutRead int) Option {
	return func(c *WebsocketConfig) {
		c.TimeoutRead = timeoutRead
	}
}

func NewWebsocketOption(options ...Option) {
	defaultWebsocketConfig = &WebsocketConfig{}
	for _, option := range options {
		option(defaultWebsocketConfig)
	}
}

var defaultWebsocketConfig *WebsocketConfig

func GetDefaultWebsocketConfig() *WebsocketConfig {
	return defaultWebsocketConfig
}
