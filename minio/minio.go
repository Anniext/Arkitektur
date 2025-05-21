package minio

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
	Token           string
}

type Option func(*OSSConfig)

func WithTokenOption(token string) Option {
	return func(c *OSSConfig) {
		c.Token = token
	}
}

func WithEndpointOption(endpoint string) Option {
	return func(c *OSSConfig) {
		c.Endpoint = endpoint
	}
}

func WithAccessKeyIDOption(accessKeyID string) Option {
	return func(c *OSSConfig) {
		c.AccessKeyID = accessKeyID
	}
}

func WithSecretAccessKeyOption(secretAccessKey string) Option {
	return func(c *OSSConfig) {
		c.SecretAccessKey = secretAccessKey
	}
}

func WithBucketNameOption(bucketName string) Option {
	return func(c *OSSConfig) {
		c.BucketName = bucketName
	}
}

func WithUseSSLOption(useSSL bool) Option {
	return func(c *OSSConfig) {
		c.UseSSL = useSSL
	}
}

func NewOSSOption(options ...Option) {
	defaultOSSConfig = &OSSConfig{}
	for _, option := range options {
		option(defaultOSSConfig)
	}
}

var defaultOSSConfig *OSSConfig

func GetDefaultOSS() *OSSConfig {
	return defaultOSSConfig
}
