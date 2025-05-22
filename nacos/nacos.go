package nacos

type NacosConfig struct {
	Name      string
	Namespace string
	Timeout   int
	Tmp       string
	Level     string
	Scheme    string
	Path      string
	Host      string
	Port      int
	Group     string
}
type Option func(*NacosConfig)

func WithNameSpaceOption(namespace string) Option {
	return func(c *NacosConfig) {
		c.Namespace = namespace
	}
}

func WithTimeoutOption(timeout int) Option {
	return func(c *NacosConfig) {
		c.Timeout = timeout
	}
}

func WithTmpOption(tmp string) Option {
	return func(c *NacosConfig) {
		c.Tmp = tmp
	}
}

func WithLevelOption(level string) Option {
	return func(c *NacosConfig) {
		c.Level = level
	}
}

func WithSchemeOption(scheme string) Option {
	return func(c *NacosConfig) {
		c.Scheme = scheme
	}
}

func WithPathOption(path string) Option {
	return func(c *NacosConfig) {
		c.Path = path
	}
}

func WithHostOption(host string) Option {
	return func(c *NacosConfig) {
		c.Host = host
	}
}

func WithPortOption(port int) Option {
	return func(c *NacosConfig) {
		c.Port = port
	}
}

func WithGroupOption(group string) Option {
	return func(c *NacosConfig) {
		c.Group = group
	}
}
func WithNameOption(name string) Option {
	return func(c *NacosConfig) {
		c.Name = name
	}
}

func NewNacosOption(options ...Option) {
	defaultNacosConfig = &NacosConfig{}
	for _, option := range options {
		option(defaultNacosConfig)
	}
}

var defaultNacosConfig *NacosConfig

func GetDefaultNacosConfig() *NacosConfig {
	return defaultNacosConfig
}
