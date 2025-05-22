package nacos

import (
	"context"
	"github.com/Anniext/Arkitektur/system/config"
	"github.com/Anniext/Arkitektur/system/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"sync"
)

func GetClientConfig() *constant.ClientConfig {
	nacosCnf := GetDefaultNacosConfig()
	return constant.NewClientConfig(
		constant.WithNamespaceId(nacosCnf.Namespace),
		constant.WithTimeoutMs(uint64(nacosCnf.Timeout)),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(nacosCnf.Tmp+"/tmp/nacos/log"),
		constant.WithCacheDir(nacosCnf.Tmp+"/tmp/nacos/cache"),
		constant.WithLogLevel(nacosCnf.Level),
	)
}

func GetServerConfig() []constant.ServerConfig {
	nacosCnf := GetDefaultNacosConfig()
	return []constant.ServerConfig{
		*constant.NewServerConfig(
			nacosCnf.Host,
			uint64(nacosCnf.Port),
			constant.WithScheme(nacosCnf.Scheme),
			constant.WithContextPath(nacosCnf.Path),
		),
	}
}

func InitDefaultNacos() error {
	ctx := context.Background()
	nacosCnf := GetDefaultNacosConfig()
	nacosDataId := nacosCnf.Name

	namingConfig, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  GetClientConfig(),
			ServerConfigs: GetServerConfig(),
		})
	if err != nil {
		return err
	}

	var serverConfig string
	serverConfig, err = namingConfig.GetConfig(vo.ConfigParam{
		DataId: nacosDataId + ".yaml",
		Group:  nacosCnf.Group,
	})

	go func(ctx context.Context) {
		var mutex sync.Mutex
		err = namingConfig.ListenConfig(vo.ConfigParam{
			DataId: nacosDataId + ".yaml",
			Group:  nacosCnf.Group,
			OnChange: func(namespace, group, dataId, data string) {
				err = yaml.Unmarshal([]byte(data), config.GetServerConfig())
				if err != nil {
					log.Error("Failed to unmarshal config:", err)
					return
				}

				mutex.Lock()
				defer mutex.Unlock()
				// 刷新数据库配置
				log.Info("switch config is successful")
			},
		})
		if err != nil {
			log.Error("failed to listenConfig ")
			return
		}
		<-ctx.Done()
	}(ctx)

	err = yaml.Unmarshal([]byte(serverConfig), config.GetServerConfig())
	if err != nil {
		panic(err)
	}
	return nil
}
