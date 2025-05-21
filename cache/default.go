package cache

import (
	"context"
	"fmt"
	"github.com/Anniext/Arkitektur/system/log"
	"github.com/redis/go-redis/v9"
	"time"
)

var defaultRedis *redis.Client

func InitDefaultRedis() error {
	cacheConfig := GetDefaultCache()
	defaultRedis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password: defaultCache.Password,
		DB:       defaultCache.DB,
		PoolSize: defaultCache.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := defaultRedis.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis连接失败: %s", err)
	}

	log.Infoln("redis server is running")
	return nil
}

func GetDefaultRedis() *redis.Client {
	return defaultRedis
}
