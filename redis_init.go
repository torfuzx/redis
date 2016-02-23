package redis

import (
	"time"

	"hotpu.cn/xkefu/common/config"
	"hotpu.cn/xkefu/common/log"
)

var (
	cache *Cache
)

// InitFromConfig initialize seelog configuration by config
func InitFromConfig(server, password string, maxIdle, maxActive int64, idleTimeout time.Duration) error {
	log.Info("common.store.redis", "InitFromConfig", "redis.InitFromConfig > run...")

	var err error

	if cache == nil {
		builder := &Builder{
			Server:      server,
			Password:    password,
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
		}

		cache, err = builder.Build()
		if err != nil {
			return err
		}

		// cache.SetLogLevel("debug")
	}

	return nil
}

func Get() (*Cache, error) {
	if cache == nil {
		return nil, config.ErrRedisNotInitialized
	}

	return cache, nil
}
