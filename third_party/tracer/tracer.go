package tracer

//import (
//	"github.com/opentracing/opentracing-go"
//	jconfig "github.com/uber/jaeger-client-go/config"
//	"github.com/uber/jaeger-client-go/rpcmetrics"
//	"github.com/uber/jaeger-lib/metrics"
//)

//
//func InitJaeger() func() {
//	jaegerHost := os.Getenv("OTEL_JAEGER_ENDPOINT")
//	jaegerServiceName := os.Getenv("OTEL_JAEGER_SERVICE_NAME")
//	jaegerExporter := os.Getenv("OTEL_EXPORTER")
//
//	// Create and install Jaeger export pipeline
//	flush, err := jaeger.InstallNewPipeline(
//		jaeger.WithCollectorEndpoint(jaegerHost + "/api/traces"),
//		jaeger.WithProcess(jaeger.Process{
//			ServiceName: jaegerServiceName,
//			Tags: []label.KeyValue{
//				label.String("exporter", jaegerExporter),
//				label.Float64("float", 312.23),
//			},
//		}),
//		jaeger.WithSDK(&sdktrace.Config{
//			DefaultSampler: sdktrace.AlwaysSample(),
//		}),
//	)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return func() {
//		flush()
//	}
//}

//func InitAutoJaeger(tracerKey string) (opentracing.Tracer, error) {
//	cfg := &jconfig.Configuration{
//		ServiceName: "tracerKey",
//		Sampler: &jconfig.SamplerConfig{
//			Type:              "const",
//			Param:             1,
//			SamplingServerURL: "localhost:5778",
//		},
//		Reporter: &jconfig.ReporterConfig{
//			LogSpans:           true,
//			CollectorEndpoint:  "http://localhost:14268/api/traces",
//			LocalAgentHostPort: "localhost:6831",
//		},
//	}
//
//	var metricsFactory metrics.Factory
//	//metricsFactory.Namespace(metrics.NSOptions{
//	//	Name: tracerKey,
//	//	Tags: map[string]string{},
//	//})
//	//metricsFactory := prometheus.New()
//	tracer, _, err := cfg.NewTracer(
//		//jconfig.Logger(jaeger.StdLogger),
//		//	jconfig.Metrics(metricsFactory),
//		jconfig.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
//	)
//	return tracer, err
//}
