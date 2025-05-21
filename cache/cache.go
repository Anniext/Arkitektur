package cache

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}
type Option func(*CacheConfig)

func WithHostOption(host string) Option {
	return func(c *CacheConfig) {
		c.Host = host
	}
}

func WithPortOption(port int) Option {
	return func(c *CacheConfig) {
		c.Port = port
	}
}

func WithPasswordOption(password string) Option {
	return func(c *CacheConfig) {
		c.Password = password
	}
}

func WithDBOption(db int) Option {
	return func(c *CacheConfig) {
		c.DB = db
	}
}

func WithPoolSizeOption(pollSize int) Option {
	return func(c *CacheConfig) {
		c.PoolSize = pollSize
	}
}

func NewCacheOption(options ...Option) {
	defaultCache = &CacheConfig{}
	for _, option := range options {
		option(defaultCache)
	}
}

var defaultCache *CacheConfig

func GetDefaultCache() *CacheConfig {
	return defaultCache
}
