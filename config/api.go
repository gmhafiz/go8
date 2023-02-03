package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	Name              string        `default:"go8_api"`
	Host              string        `default:"0.0.0.0"`
	Port              string        `default:"3080"`
	ReadHeaderTimeout time.Duration `split_words:"true" default:"60s"`
	GracefulTimeout   time.Duration `split_words:"true" default:"8s"`

	RequestLog bool `split_words:"true" default:"false"`
	RunSwagger bool `split_words:"true" default:"true"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
