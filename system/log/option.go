package log

type SlogConfig struct {
	ProdLevel  string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	Mode       string
}

type Option func(*SlogConfig)

func WithProdLevelOption(prodLevel string) Option {
	return func(c *SlogConfig) {
		c.ProdLevel = prodLevel
	}
}

func WithFilenameOption(filename string) Option {
	return func(c *SlogConfig) {
		c.Filename = filename
	}
}

func WithMaxSizeOption(maxSize int) Option {
	return func(c *SlogConfig) {
		c.MaxSize = maxSize
	}
}

func WithMaxAgeOption(maxAge int) Option {
	return func(c *SlogConfig) {
		c.MaxAge = maxAge
	}
}

func WithMaxBackupsOption(maxBackups int) Option {
	return func(c *SlogConfig) {
		c.MaxBackups = maxBackups
	}
}

func WithCompressOption(compress bool) Option {
	return func(c *SlogConfig) {
		c.Compress = compress
	}
}

func WithModeOption(mode string) Option {
	return func(c *SlogConfig) {
		c.Mode = mode
	}
}

func NewSlogOption(options ...Option) *SlogConfig {
	config := &SlogConfig{}
	for _, option := range options {
		option(config)
	}
	return config
}
