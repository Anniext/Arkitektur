package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	serverConfig *ServerConfig
)

// InitSystemConfig 初始化配置文件
func InitSystemConfig(name, mode, path string) error {
	v := viper.New()

	envString := fmt.Sprintf("%s-%s", name, mode)
	pwd, _ := os.Getwd()
	configPath := filepath.Join(pwd, path)
	v.AddConfigPath(configPath)
	v.SetConfigName(envString)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	conf := &ServerConfig{}
	if err := v.Unmarshal(conf); err != nil {
		fmt.Println("Unmarshal is err")
		return err
	}

	serverConfig = conf
	return nil
}

// GetServerConfig 获取配置文件
func GetServerConfig() *ServerConfig {
	return serverConfig
}

// GetMysqlInfo 获取mysql配置文件
func GetMysqlInfo() MysqlInfo {
	return serverConfig.MysqlInfo
}

// GetRedisInfo 获取redis配置文件
func GetRedisInfo() RedisInfo {
	return serverConfig.RedisInfo
}

// GetNacosInfo 获取nacos配置文件
func GetNacosInfo() NacosInfo {
	return serverConfig.NacosInfo
}

// GetMinioInfo 获取minio配置文件
func GetMinioInfo() MinioConfig {
	return serverConfig.MinioInfo
}

// GetCasbinInfo 获取casbin配置文件
func GetCasbinInfo() CasbinInfo {
	return serverConfig.Casbin
}

// GetMqttInfo 获取mqtt配置文件
func GetMqttInfo() MqttConfig {
	return serverConfig.MqttInfo
}

// GetWebsocketInfo 获取websocket配置文件
func GetWebsocketInfo() WebsocketInfo {
	return serverConfig.WebsocketInfo
}
