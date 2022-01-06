package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Cache struct {
	Host      string `default:"0.0.0.0"`
	Port      string `default:"6379"`
	Name      int    `default:"1"`
	User      string
	Pass      string
	CacheTime int `default:"5"`
}

func NewCache() Cache {
	var cache Cache
	envconfig.MustProcess("REDIS", &cache)

	return cache
}
