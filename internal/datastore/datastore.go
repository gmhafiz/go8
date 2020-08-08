package datastore

import (
	"database/sql"
	"fmt"
)

type Database struct {
	Driver  string `yaml:"DRIVER"`
	Host    string `yaml:"HOST"`
	Port    string `yaml:"PORT"`
	Name    string `yaml:"NAME"`
	User    string `yaml:"USER"`
	Pass    string `yaml:"PASS"`
	SslMode string `yaml:"SSL_MODE"`
}

func NewService(cfg *Database) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
		cfg.Driver,
		cfg.Host,
		cfg.Name,
		cfg.SslMode,
		cfg.User,
		cfg.Pass,
	)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
