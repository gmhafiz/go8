package jobs

import (
	"github.com/gomodule/redigo/redis"

	"go8ddd/configs"
)

type RedisClient struct {
	Pool *redis.Pool
	Conn redis.Conn
}

func newPool(cfg *configs.Configs) *RedisClient {
	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":"+cfg.Cache.Port)
		},
	}

	conn, _ := redis.Dial("tcp", ":"+cfg.Cache.Port)

	return &RedisClient{
		Pool: pool,
		Conn: conn,
	}
}
