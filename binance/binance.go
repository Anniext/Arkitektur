package binance

type BinanceConfig struct {
	Proxy     string
	ApiKey    string
	SecretKey string
}
type Option func(*BinanceConfig)

func WithProxyOption(proxy string) Option {
	return func(c *BinanceConfig) {
		c.Proxy = proxy
	}
}
func WithApiKeyOption(apiKey string) Option {
	return func(c *BinanceConfig) {
		c.ApiKey = apiKey
	}
}
func WithSecretKeyOption(secretKey string) Option {
	return func(c *BinanceConfig) {
		c.SecretKey = secretKey
	}
}

func NewBinanceOption(options ...Option) {
	defaultBinanceConfig = &BinanceConfig{}
	for _, option := range options {
		option(defaultBinanceConfig)
	}
}

var defaultBinanceConfig *BinanceConfig

func GetDefaultBinanceConfig() *BinanceConfig {
	return defaultBinanceConfig
}
