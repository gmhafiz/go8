package configs

import (
	"os"
	"strconv"
	"time"
)

type Api struct {
	Name              string
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	RequestLog        bool
	RunSwagger        bool
}

func API() *Api {
	apiReadTimeout, _ := strconv.Atoi(os.Getenv("API_READ_TIMEOUT"))
	apiReadHeaderTimeout, _ := strconv.Atoi(os.Getenv("API_READ_HEADER_TIMEOUT"))
	apiWriteTimeout, _ := strconv.Atoi(os.Getenv("API_WRITE_TIMEOUT"))
	apiIdleTimeout, _ := strconv.Atoi(os.Getenv("API_IDLE_TIMEOUT"))
	requestLog, _ := strconv.ParseBool(os.Getenv("API_REQUEST_LOG"))
	runSwagger, _ := strconv.ParseBool(os.Getenv("RUN_SWAGGER"))

	return &Api{
		Name:              os.Getenv("API_HOST"),
		Host:              os.Getenv("API_HOST"),
		Port:              os.Getenv("API_PORT"),
		ReadTimeout:       time.Duration(apiReadTimeout),
		ReadHeaderTimeout: time.Duration(apiReadHeaderTimeout),
		WriteTimeout:      time.Duration(apiWriteTimeout),
		IdleTimeout:       time.Duration(apiIdleTimeout),
		RequestLog:        requestLog,
		RunSwagger:        runSwagger,
	}
}
