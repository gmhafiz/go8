package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"go8ddd/configs"
)

func New(cfg *configs.Configs) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		//Network:            "",
		Addr:      fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Dialer:    nil,
		OnConnect: nil,
		Username:  cfg.Cache.User,
		Password:  cfg.Cache.Pass,
		DB:        cfg.Cache.Name,
		//MaxRetries:         0,
		//MinRetryBackoff:    0,
		//MaxRetryBackoff:    0,
		//DialTimeout:        0,
		//ReadTimeout:        0,
		//WriteTimeout:       0,
		//PoolSize:           0,
		//MinIdleConns:       0,
		//MaxConnAge:         0,
		//PoolTimeout:        0,
		//IdleTimeout:        0,
		//IdleCheckFrequency: 0,
		//TLSConfig:          nil,
		//Limiter:            nil,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
