package configs

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

type Cache struct {
	Host      net.IP `default:"0.0.0.0"`
	Port      string `default:"6379"`
	Name      int
	User      string
	Pass      string
	CacheTime int
}

func NewCache() Cache {
	var cache Cache
	envconfig.MustProcess("REDIS", &cache)

	return cache
}
