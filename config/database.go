package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver            string
	Host              string
	Port              string
	Name              string
	User              string
	Pass              string
	SslMode           string `default:"disable"`
	MaxConnectionPool int    `default:"4"`
}

func DataStore() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
