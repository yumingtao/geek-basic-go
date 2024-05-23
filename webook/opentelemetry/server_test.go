package opentelemetry

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)
	prop := newPropagator()
	// 在客户端和服务端之间传递tracing的相关信息
	otel.SetTextMapPropagator(prop)
	// 初始化trace provider，用来在打点时构建trace
	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	defer func(tp *trace.TracerProvider, ctx context.Context) {
		_ = tp.Shutdown(ctx)
	}(tp, context.Background())
	otel.SetTracerProvider(tp)

	server := gin.Default()
	server.GET("/test", func(ginCtx *gin.Context) {
		tracer := otel.Tracer("yumingtao/opentelemetry")
		var ctx context.Context = ginCtx
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()
		time.Sleep(time.Second)
		span.AddEvent("发生了某事")
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()
		subSpan.SetAttributes(attribute.String("attr1", "value1"))
		time.Sleep(time.Millisecond * 300)
		ginCtx.String(http.StatusOK, "测试 span")
	})
	err = server.Run(":8082")
	if err != nil {
		panic(err)
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
