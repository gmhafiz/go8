package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Cache struct {
	Enable    bool   `default:"false"`
	Host      string `default:"0.0.0.0"`
	Port      string `default:"6379"`
	Hosts     []string
	Name      int `default:"1"`
	User      string
	Pass      string
	CacheTime time.Duration `default:"5s"`
}

func NewCache() Cache {
	var cache Cache
	envconfig.MustProcess("REDIS", &cache)

	if strings.Contains(cache.Host, ",") {
		cache.Hosts = strings.Split(cache.Host, ",")
	}

	return cache
}
