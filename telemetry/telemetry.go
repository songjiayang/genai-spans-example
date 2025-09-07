package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	trace2 "go.opentelemetry.io/otel/trace"
)

// InitTracer 初始化OpenTelemetry tracer
func InitTracer() func() {
	// 创建console exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("Failed to create stdout exporter: %v", err)
	}

	// 创建trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	// 设置全局trace provider
	otel.SetTracerProvider(tp)

	// 返回cleanup函数
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

// GetTracer 获取tracer实例
func GetTracer(name string) trace2.Tracer {
	return otel.Tracer(name)
}
