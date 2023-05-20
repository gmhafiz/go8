package config

import "github.com/kelseyhightower/envconfig"

type Cors struct {
	AllowedOrigins []string `split_words:"true"`
}

func NewCors() Cors {
	var c Cors
	envconfig.MustProcess("CORS", &c)

	return c
}
