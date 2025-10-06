package telemetry

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	trace2 "go.opentelemetry.io/otel/trace"
)

// ExporterType 定义导出器类型
type ExporterType string

const (
	ExporterConsole ExporterType = "console"
	ExporterHTTP    ExporterType = "http"
	ExporterAuto    ExporterType = "auto"
)

// Config 定义telemetry配置
type Config struct {
	ExporterType ExporterType
	Endpoint     string
	ServiceName  string
}

// InitTracer 初始化OpenTelemetry tracer
func InitTracer() func() {
	return InitTracerWithConfig(Config{
		ExporterType: ExporterConsole,
		ServiceName:  "gen-ai-example",
	})
}

// InitTracerWithConfig 使用配置初始化OpenTelemetry tracer
func InitTracerWithConfig(config Config) func() {
	if config.ServiceName == "" {
		config.ServiceName = "gen-ai-example"
	}

	var exporter trace.SpanExporter
	var err error

	// 根据配置创建导出器
	switch config.ExporterType {
	case ExporterHTTP:
		if config.Endpoint == "" {
			config.Endpoint = "http://localhost:4318"
		}
		exporter, err = createHTTPOutputExporter(config.Endpoint)
	case ExporterAuto:
		if config.Endpoint == "" {
			config.Endpoint = "http://localhost:4318"
		}
		// 优先尝试HTTP，失败后回退到console
		if exporter, err = createHTTPOutputExporter(config.Endpoint); err != nil {
			log.Printf("Failed to create HTTP exporter, falling back to console: %v", err)
			exporter, err = createConsoleExporter()
		}
	case ExporterConsole, "":
		exporter, err = createConsoleExporter()
	default:
		log.Fatalf("Unsupported exporter type: %s", config.ExporterType)
	}

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	// 创建trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
		)),
	)

	// 设置全局trace provider
	otel.SetTracerProvider(tp)

	// 返回cleanup函数
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		if err := exporter.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down exporter: %v", err)
		}
	}
}

// createConsoleExporter 创建console导出器
func createConsoleExporter() (*stdouttrace.Exporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
}

// createHTTPOutputExporter 创建HTTP导出器
func createHTTPOutputExporter(endpoint string) (trace.SpanExporter, error) {
	// 验证endpoint URL
	if _, err := url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %v", err)
	}

	// 确保endpoint不包含路径和查询参数
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint URL: %v", err)
	}

	cleanEndpoint := parsedURL.Hostname()
	if parsedURL.Port() != "" {
		cleanEndpoint = cleanEndpoint + ":" + parsedURL.Port()
	}

	return otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cleanEndpoint),
		otlptracehttp.WithInsecure(),
	)
}

// GetTracer 获取tracer实例
func GetTracer(name string) trace2.Tracer {
	return otel.Tracer(name)
}

// GetConfigFromEnv 从环境变量获取配置
func GetConfigFromEnv() Config {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	serviceName := os.Getenv("OTEL_SERVICE_NAME")

	var exporterType ExporterType
	switch os.Getenv("OTEL_TRACES_EXPORTER") {
	case "otlp", "http":
		exporterType = ExporterHTTP
	case "console", "":
		exporterType = ExporterConsole
	default:
		exporterType = ExporterType(os.Getenv("OTEL_TRACES_EXPORTER"))
	}

	if endpoint != "" && exporterType != ExporterConsole {
		// 如果设置了endpoint，默认使用HTTP导出器
		exporterType = ExporterHTTP
	}

	return Config{
		ExporterType: exporterType,
		Endpoint:     endpoint,
		ServiceName:  serviceName,
	}
}
