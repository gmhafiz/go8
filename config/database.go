package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver                 string
	Host                   string
	Port                   string
	Name                   string
	User                   string
	Pass                   string
	SslMode                string        `default:"disable"`
	MaxConnectionPool      int           `default:"4"`
	MaxIdleConnections     int           `default:"4"`
	ConnectionsMaxLifeTime time.Duration `default:"300s"`
}

func DataStore() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
