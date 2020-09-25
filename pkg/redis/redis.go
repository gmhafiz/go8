package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `yaml:"HOST"`
	Port     string `yaml:"PORT"`
	Name     string `yaml:"NAME"`
	Username string `yaml:"USER"`
	Password string `yaml:"PASS"`
}

func NewClient(cfg *Config) (*redis.Client, error) {
	var ctx = context.Background()

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		//Network:            "",
		Addr:      address,
		Dialer:    nil,
		OnConnect: nil,
		Username:  cfg.Username,
		Password:  cfg.Password,
		//DB:        0,
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

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	// todo: return struct
	return rdb, nil
}
