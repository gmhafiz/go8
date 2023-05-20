package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver                 string        `required:"true"`
	Host                   string        `default:"localhost"`
	Port                   uint16        `default:"5432"`
	Name                   string        `default:"postgres"`
	TestName               string        `split_words:"true" default:"test"`
	User                   string        `default:"postgres"`
	Pass                   string        `default:"password"`
	SslMode                string        `split_words:"true" default:"disable"`
	MaxConnectionPool      int           `split_words:"true" default:"4"`
	MaxIdleConnections     int           `split_words:"true" default:"4"`
	ConnectionsMaxLifeTime time.Duration `split_words:"true" default:"300s"`
}

func DataStore() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
