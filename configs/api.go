package configs

import (
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	Name              string        `default:"go8"`
	Host              net.IP        `default:"0.0.0.0"`
	Port              string        `default:"3080"`
	ReadTimeout       time.Duration `default:"5s"`
	ReadHeaderTimeout time.Duration `default:"5s"`
	WriteTimeout      time.Duration `default:"10s"`
	IdleTimeout       time.Duration `default:"120s"`
	RequestLog        bool          `default:"false"`
	RunSwagger        bool          `default:"true"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
