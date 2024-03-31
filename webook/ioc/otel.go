package ioc

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	trace2 "go.opentelemetry.io/otel/trace/noop"
	"time"
)

// InitOTEL 返回一个关闭函数，让调用者关闭的时候来决定这个ctx
func InitOTEL() func(ctx context.Context) {
	res, err := newResource("webook", "v0.0.1")
	if err != nil {
		panic(err)
	}
	prop := newPropagator()
	// 在客户端和服务端之间传递tracing的相关信息
	otel.SetTextMapPropagator(prop)
	// 初始化trace provider，用来在打点时构建trace
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		_ = tp.Shutdown(ctx)
	}
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
func newResource(serviceName string, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion)),
	)
}

func NewTracer() trace2.Tracer {
	return trace2.Tracer{}
}
