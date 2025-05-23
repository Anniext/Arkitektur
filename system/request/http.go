package request

type HttpConfig struct {
	ApiUrl string
	Host   string
	Key    string
}

type Option func(*HttpConfig)

func WithApiUrlOption(apiUrl string) Option {
	return func(c *HttpConfig) {
		c.ApiUrl = apiUrl
	}
}

func WithKeyOption(key string) Option {
	return func(c *HttpConfig) {
		c.Key = key
	}
}

func WithHostOption(host string) Option {
	return func(c *HttpConfig) {
		c.Host = host
	}
}

func NewHttpOption(options ...Option) {
	defaultHttpConfig = &HttpConfig{}
	for _, option := range options {
		option(defaultHttpConfig)
	}
}

var defaultHttpConfig *HttpConfig

func GetDefaultHttpConfig() *HttpConfig {
	return defaultHttpConfig
}
