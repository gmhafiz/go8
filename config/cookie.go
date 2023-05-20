package config

import (
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Session struct {
	Name     string          `split_words:"true" default:"session"`
	Path     string          `default:"/"`
	Domain   string          `default:""`
	Secret   string          `required:"false"`
	Duration time.Duration   `default:"24h"`
	HttpOnly bool            `split_words:"true" default:"true"`
	Secure   bool            `default:"true"`
	SameSite SameSiteDecoder `split_words:"true" default:"lax"`
}

func NewSession() Session {
	var a Session
	envconfig.MustProcess("SESSION", &a)

	return a
}

type SameSiteDecoder http.SameSite

func (sd *SameSiteDecoder) Decode(value string) error {
	switch value {
	case "default":
		*sd = SameSiteDecoder(http.SameSiteDefaultMode)
	case "lax":
		*sd = SameSiteDecoder(http.SameSiteLaxMode)
	case "strict":
		*sd = SameSiteDecoder(http.SameSiteStrictMode)
	case "none":
	case "":
		*sd = SameSiteDecoder(http.SameSiteLaxMode)
	}

	return nil
}
