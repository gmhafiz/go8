package logger

import (
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func New(version string) zerolog.Logger {
	logger := httplog.NewLogger("go8", httplog.Options{
		JSON:    false, // switch to false for a human readable log format
		Concise: true,
		Tags:    map[string]string{"version": version},
	})
	return logger.With().Caller().Logger()
}
