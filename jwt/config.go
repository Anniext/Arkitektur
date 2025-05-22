package jwt

type JwtConfig struct {
	JwtSigningKey string
}
type Option func(*JwtConfig)

func WithJwtSigningKeyOption(jwtSigningKey string) Option {
	return func(c *JwtConfig) {
		c.JwtSigningKey = jwtSigningKey
	}
}

func NewCacheOption(options ...Option) {
	defaultJwtConfig = &JwtConfig{}
	for _, option := range options {
		option(defaultJwtConfig)
	}
}

var defaultJwtConfig *JwtConfig

func GetDefaultJwtConfig() *JwtConfig {
	return defaultJwtConfig
}
