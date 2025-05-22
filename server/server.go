package server

type GinConfig struct {
	Addr string
	Port int
}
type Option func(*GinConfig)

func WithAddrOption(addr string) Option {
	return func(c *GinConfig) {
		c.Addr = addr
	}
}

func WithPortOption(port int) Option {
	return func(c *GinConfig) {
		c.Port = port
	}
}

func NewGinOption(options ...Option) {
	defaultGinConfig = &GinConfig{}
	for _, option := range options {
		option(defaultGinConfig)
	}
}

var defaultGinConfig *GinConfig

func GetDefaultGinConfig() *GinConfig {
	return defaultGinConfig
}
