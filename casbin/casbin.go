package casbin

type CasbinConfig struct {
	ModelPath string
}
type Option func(*CasbinConfig)

func WithModelPathOption(modelPath string) Option {
	return func(c *CasbinConfig) {
		c.ModelPath = modelPath
	}
}

func NewCacheOption(options ...Option) {
	defaultCasbinConfig = &CasbinConfig{}
	for _, option := range options {
		option(defaultCasbinConfig)
	}
}

var defaultCasbinConfig *CasbinConfig

func GetDefaultCasbinConfig() *CasbinConfig {
	return defaultCasbinConfig
}
