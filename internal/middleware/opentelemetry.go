package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Cancel func()
}

func Otlp(enable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)

				if !enable {
					return
				}

				defaultCtx := baggage.ContextWithoutBaggage(r.Context())
				routePattern := chi.RouteContext(defaultCtx).RoutePattern()
				span := trace.SpanFromContext(defaultCtx)
				span.SetName(routePattern)
				span.SetAttributes(semconv.HTTPTarget(r.URL.String()), semconv.HTTPRoute(routePattern))
				labeler, ok := otelhttp.LabelerFromContext(defaultCtx)
				if ok {
					labeler.Add(semconv.HTTPRoute(routePattern))
				}
			}),
			"",
		)
	}
}
