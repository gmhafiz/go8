package configs

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultAPIHost              = "0.0.0.0"
	defaultAPIPort              = 3080
	defaultApiReadTimeout       = 5
	defaultApiReadHeaderTimeout = 5
	defaultApiWriteTimeout      = 10
	defaultApiIdleTimeout       = 120
	defaultRequestLog           = false
	defaultRunSwagger           = true
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
	apiName := os.Getenv("API_NAME")
	if apiName == "" {
		apiName = "Go API"
	}

	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = defaultAPIHost
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = strconv.Itoa(defaultAPIPort)
	}

	apiReadTimeout, err := strconv.Atoi(os.Getenv("API_READ_TIMEOUT"))
	if err != nil {
		apiReadTimeout = defaultApiReadTimeout
	}

	apiReadHeaderTimeout, err := strconv.Atoi(os.Getenv("API_READ_HEADER_TIMEOUT"))
	if err != nil {
		apiReadHeaderTimeout = defaultApiReadHeaderTimeout
	}

	apiWriteTimeout, err := strconv.Atoi(os.Getenv("API_WRITE_TIMEOUT"))
	if err != nil {
		apiWriteTimeout = defaultApiWriteTimeout
	}

	apiIdleTimeout, err := strconv.Atoi(os.Getenv("API_IDLE_TIMEOUT"))
	if err != nil {
		apiIdleTimeout = defaultApiIdleTimeout
	}

	requestLog, err := strconv.ParseBool(os.Getenv("API_REQUEST_LOG"))
	if err != nil {
		requestLog = defaultRequestLog
	}

	runSwagger, err := strconv.ParseBool(os.Getenv("RUN_SWAGGER"))
	if err != nil {
		runSwagger = defaultRunSwagger
	}

	return &Api{
		Name:              apiName,
		Host:              apiHost,
		Port:              apiPort,
		ReadTimeout:       time.Duration(apiReadTimeout),
		ReadHeaderTimeout: time.Duration(apiReadHeaderTimeout),
		WriteTimeout:      time.Duration(apiWriteTimeout),
		IdleTimeout:       time.Duration(apiIdleTimeout),
		RequestLog:        requestLog,
		RunSwagger:        runSwagger,
	}
}
