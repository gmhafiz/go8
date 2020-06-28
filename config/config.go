package config

import (
	"log"
	"os"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Conf struct {
	Debug   bool `env:"DEBUG,required"`
	Testing bool `env:"TESTING,false"`
	Server  serverConf
	Db      dbConf
}

type serverConf struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
	SslMode  string `env:"SSL_Mode,default=disable"`

	TestHost     string `env:"TEST_DB_HOST,required"`
	TestPort     int    `env:"TEST_DB_PORT,required"`
	TestUsername string `env:"TEST_DB_USER,required"`
	TestPassword string `env:"TEST_DB_PASS,required"`
	TestDbName   string `env:"TEST_DB_NAME,required"`
	TestSslMode  string `env:"TEST_SSL_Mode,default=disable"`
}

func AppConfig() *Conf {
	cwd, err := os.Getwd()
	log.Println(cwd)
	if err != nil {
		log.Fatalln(err)
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config Conf
	if err := envdecode.StrictDecode(&config); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &config
}
