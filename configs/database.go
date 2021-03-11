package configs

import (
	"log"
	"os"
	"strconv"
)

type Database struct {
	Driver      string
	Host        string
	Port        string
	Name        string
	User        string
	Pass        string
	SslMode     string
	MaxConnPool int
}

func DataStore() *Database {
	num, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTION_POOL"))
	if err != nil {
		log.Fatal(err)
	}
	return &Database{
		Driver:      os.Getenv("DB_DRIVER"),
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Name:        os.Getenv("DB_NAME"),
		User:        os.Getenv("DB_USER"),
		Pass:        os.Getenv("DB_PASS"),
		SslMode:     os.Getenv("DB_SSL_MODE"),
		MaxConnPool: num,
	}
}
