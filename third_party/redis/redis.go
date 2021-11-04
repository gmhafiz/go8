package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/gmhafiz/go8/configs"
)

func New(cfg configs.Cache) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.Name,
	})
}
