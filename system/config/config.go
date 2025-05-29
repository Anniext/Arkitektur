package config

type ServerConfig struct {
	Name          string        `mapstructure:"name" json:"name" yaml:"name"`
	JwtSigningKey string        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	BarkUrl       string        `mapstructure:"bark" json:"bark" yaml:"bark"`
	ZoneId        int32         `mapstructure:"zone_id" json:"zone_id" yaml:"zone_id"`
	MysqlInfo     MysqlInfo     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	RedisInfo     RedisInfo     `mapstructure:"redis" json:"redis" yaml:"redis"`
	NacosInfo     NacosInfo     `mapstructure:"nacos" json:"nacos" yaml:"nacos"`
	MinioInfo     MinioConfig   `mapstructure:"oss" json:"oss" yaml:"oss"`
	Casbin        CasbinInfo    `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	MqttInfo      MqttConfig    `mapstructure:"mqtt" json:"mqtt" yaml:"mqtt"`
	WebsocketInfo WebsocketInfo `mapstructure:"websocket" json:"websocket" yaml:"websocket"`
	TweetInfo     TweetInfo     `mapstructure:"tweet" json:"tweet" yaml:"tweet"`
	BinanceInfo   BinanceInfo   `mapstructure:"binance" json:"binance" yaml:"binance"`
}

type TweetInfo struct {
	MemoryUrl string `mapstructure:"memory_url" json:"memory_url" yaml:"memory_url"`
	ApiUrl    string `mapstructure:"api_url" json:"api_url" yaml:"api_url"`
	Key       string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
	Host      string `mapstructure:"api_host" json:"api_host" yaml:"api_host"`
}

// MysqlInfo mysql配置文件
type MysqlInfo struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       string `mapstructure:"db" json:"db" yaml:"db"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// RedisInfo redis配置文件
type RedisInfo struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// NacosInfo nacos配置文件 可选
type NacosInfo struct {
	Enable      bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`
	Namespace   string `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	Group       string `mapstructure:"group" json:"group" yaml:"group"`
	DataId      string `mapstructure:"dataId" json:"dataId" yaml:"dataId"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId" yaml:"namespaceId"`
}

// MinioConfig minio配置文件 可选
type MinioConfig struct {
	Enable          bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key" json:"secret-access-key" yaml:"secret-access-key"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	UseSSL          bool   `mapstructure:"use-ssl" json:"use-ssl" yaml:"use-ssl"`
	Token           string `mapstructure:"token" json:"token" yaml:"token"`
}

type CasbinInfo struct {
	Enable    bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	ModelPath string `mapstructure:"modelPath" json:"modelPath" yaml:"modelPath"`
}

type MqttConfig struct {
	ServerID             string `mapstructure:"server_id" json:"server_id" yaml:"server_id"`                                     // mqtt client id
	BrokerURL            string `mapstructure:"broker_url" json:"broker_url" yaml:"broker_url"`                                  // mqtt服务器
	KeepAlive            int64  `mapstructure:"keep_alive" json:"keep_alive" yaml:"keep_alive"`                                  // 心跳包
	QoS                  byte   `mapstructure:"qos" json:"qos" yaml:"qos"`                                                       // 消息质量
	Retain               bool   `mapstructure:"retain" json:"retain" yaml:"retain"`                                              // 保留消息标志
	AutoReconnect        bool   `mapstructure:"auto_reconnect" json:"auto_reconnect" yaml:"auto_reconnect"`                      // 是否自重连
	ConnectTimeout       int64  `mapstructure:"connect_timeout" json:"connect_timeout" yaml:"connect_timeout"`                   // 连接超时时间
	MaxReconnectInterval int64  `mapstructure:"max_reconnectInterval" json:"max_reconnectInterval" yaml:"max_reconnectInterval"` // 最大重连间隔
}

type WebsocketInfo struct {
	Enable      bool `mapstructure:"enable" json:"enable" yaml:"enable"`
	Port        int  `mapstructure:"port" json:"port" yaml:"port"`
	TimeoutRead int  `mapstructure:"timeout_read" json:"timeout_read" yaml:"timeout_read"`
}
type BinanceInfo struct {
	ApiKey    int `mapstructure:"apiKey" json:"apiKey" yaml:"apiKey"`
	SecretKey int `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
}
