package config

import (
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

type OpenTelemetry struct {
	Enable             bool    `default:"false"`
	OtlpEndpoint       string  `split_words:"true" default:"0.0.0.0:4317"`
	OtlpServiceName    string  `split_words:"true" default:"go8"`
	OtlpServiceVersion string  `split_words:"true" default:"0.1.0"`
	OtlpMeterName      string  `split_words:"true" default:"go8-meter"`
	OtlpSamplerRatio   float64 `split_words:"true" default:"0.1"`
}

func NewOpenTelemetry() OpenTelemetry {
	var otel OpenTelemetry
	envconfig.MustProcess("OTEL", &otel)

	if otel.OtlpSamplerRatio < 0 || otel.OtlpSamplerRatio > 1 {
		slog.Error("trace sample ratio must be between 0 and 1")
	}

	return otel
}
