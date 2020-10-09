package configs

import (
	"os"
	"strconv"
	"time"
)

type Api struct {
	Host                 string
	Port                 string
	ApiReadTimeout       time.Duration
	ApiReadHeaderTimeout time.Duration
	ApiWriteTimeout      time.Duration
	ApiIdleTimeout       time.Duration
}

func API() *Api {
	apiReadTimeout, _ := strconv.Atoi(os.Getenv("API_READ_TIMEOUT"))
	apiReadHeaderTimeout, _ := strconv.Atoi(os.Getenv("API_READ_HEADER_TIMEOUT"))
	apiWriteTimeout, _ := strconv.Atoi(os.Getenv("API_WRITE_TIMEOUT"))
	apiIdleTimeout, _ := strconv.Atoi(os.Getenv("API_IDLE_TIMEOUT"))

	return &Api{
		Host:                 os.Getenv("API_HOST"),
		Port:                 os.Getenv("API_PORT"),
		ApiReadTimeout:       time.Duration(apiReadTimeout),
		ApiReadHeaderTimeout: time.Duration(apiReadHeaderTimeout),
		ApiWriteTimeout:      time.Duration(apiWriteTimeout),
		ApiIdleTimeout:       time.Duration(apiIdleTimeout),
	}
}
