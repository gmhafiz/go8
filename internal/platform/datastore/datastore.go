package datastore

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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

func NewService(cfg *Database) (*pgxpool.Pool, *sql.DB, error) {
	//poolcfg, err := pgxpool.ParseConfig(cfg.ConnURL())
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//poolcfg.MaxConnLifetime = cfg.IdleTimeout
	//poolcfg.MaxConns = int32(cfg.ConnPoolSize)
	//
	//dialer := &net.Dialer{KeepAlive: cfg.DialTimeout}
	//dialer.Timeout = cfg.DialTimeout
	//poolcfg.ConnConfig.DialFunc = dialer.DialContext
	//
	//pool, err := pgxpool.ConnectConfig(context.Background(), poolcfg)
	//if err != nil {
	//	return nil, nil, err
	//}

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
		return nil, nil, err
	}

	return nil, db, nil
}
