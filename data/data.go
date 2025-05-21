package data

type DBConfig struct {
	Mode   string
	Dns    string
	Driver string
}
type Option func(*DBConfig)

func WithDriverOption(driver string) Option {
	return func(c *DBConfig) {
		c.Driver = driver
	}
}

func WithDnsOption(dns string) Option {
	return func(c *DBConfig) {
		c.Dns = dns
	}
}

func WithModeOption(mode string) Option {
	return func(c *DBConfig) {
		c.Mode = mode
	}
}

func NewDBOption(options ...Option) {
	defaultDB = &DBConfig{}
	for _, option := range options {
		option(defaultDB)
	}
}

var defaultDB *DBConfig

func GetDefaultDB() *DBConfig {
	return defaultDB
}
