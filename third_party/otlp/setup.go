package otlp

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gmhafiz/go8/config"
)

// SetupOTLPExporter bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
// Reference: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/dice/otel.go
func SetupOTLPExporter(ctx context.Context, cfg config.OpenTelemetry) func() {
	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.OtlpServiceName),
			semconv.ServiceVersionKey.String(cfg.OtlpServiceVersion),
		),
	)
	if err != nil {
		log.Println(fmt.Errorf("creating resource, %v", err))
	}

	conn, err := dialGrpc(ctx, cfg)
	if err != nil {
		log.Println(fmt.Errorf("connecting to otel-collecter: %w", err))
	}

	tracerProvider, err := newTraceProvider(ctx, res, conn, cfg)
	if err != nil {
		log.Println(fmt.Errorf("creating trace provider, %v", err))
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx, res, conn)
	if err != nil {
		log.Println(fmt.Errorf("creating meter provider, %v", err))
	}

	otel.SetMeterProvider(meterProvider)

	log.Println("otlp connected.")

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
		// pushes any last exports to the receiver
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func dialGrpc(ctx context.Context, cfg config.OpenTelemetry) (*grpc.ClientConn, error) {
	log.Printf("dialing %s\n", cfg.OtlpEndpoint)
	for {
		conn, err := grpc.DialContext(ctx, cfg.OtlpEndpoint,
			// Note the use of insecure transport here. TLS is recommended in production.
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			log.Printf("gRPC connected.\n")
			return conn, nil
		}

		base, capacity := time.Second, time.Minute
		for backoff := base; err != nil; backoff <<= 1 {
			if backoff > capacity {
				backoff = capacity
			}
			jitter := rand.Int63n(int64(backoff * 3))
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			log.Println("retrying to connect to gRPC...")
			conn, err := grpc.DialContext(ctx, cfg.OtlpEndpoint,
				// Note the use of insecure transport here. TLS is recommended in production.
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err == nil {
				log.Printf("gRPC connected.\n")
				return conn, nil
			}
		}
	}
}

func newTraceProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn, cfg config.OpenTelemetry) (*trace.TracerProvider, error) {
	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := trace.NewBatchSpanProcessor(traceExp)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(cfg.OtlpSamplerRatio)),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(
			metricExporter,
			metric.WithInterval(5*time.Second),
		)),
	)
	return meterProvider, nil
}
