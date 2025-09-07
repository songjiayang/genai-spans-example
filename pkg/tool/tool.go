package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gen-ai-example/telemetry"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

type WeatherTool struct{}

func (w *WeatherTool) Name() string {
	return "get_weather"
}

func (w *WeatherTool) Description() string {
	return "Get weather information for a specified city"
}

func (w *WeatherTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	city, ok := params["city"].(string)
	if !ok {
		return nil, fmt.Errorf("missing city parameter")
	}

	time.Sleep(50 * time.Millisecond)

	return map[string]interface{}{
		"city":        city,
		"temperature": "22°C",
		"condition":   "晴天",
		"humidity":    "65%",
	}, nil
}

type CalculatorTool struct{}

func (c *CalculatorTool) Name() string {
	return "calculator"
}

func (c *CalculatorTool) Description() string {
	return "Perform basic mathematical calculations"
}

func (c *CalculatorTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("missing operation parameter")
	}

	a, ok1 := params["a"].(float64)
	b, ok2 := params["b"].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("missing or invalid numeric parameters")
	}

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	return map[string]interface{}{
		"operation": operation,
		"a":         a,
		"b":         b,
		"result":    result,
	}, nil
}

type ToolService struct {
	tools  map[string]Tool
	tracer trace.Tracer
}

func NewToolService() *ToolService {
	ts := &ToolService{
		tools:  make(map[string]Tool),
		tracer: telemetry.GetTracer("tool-service"),
	}

	ts.RegisterTool(&WeatherTool{})
	ts.RegisterTool(&CalculatorTool{})

	return ts
}

func (ts *ToolService) RegisterTool(tool Tool) {
	ts.tools[tool.Name()] = tool
}

func (ts *ToolService) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	ctx, span := ts.tracer.Start(ctx, "tool.execute",
		trace.WithAttributes(
			semconv.GenAIProviderNameOpenAI,
			semconv.GenAIOperationNameExecuteTool,
			semconv.GenAIToolName(toolName),
			semconv.GenAIToolCallID(uuid.NewString()),
		),
	)
	defer span.End()

	tool, exists := ts.tools[toolName]
	if !exists {
		span.RecordError(fmt.Errorf("tool not found: %s", toolName))
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	paramsJSON, _ := json.Marshal(params)
	span.SetAttributes(
		semconv.GenAIToolDescription(tool.Description()),
		semconv.GenAIToolType("function"),
		attribute.String("gen_ai.tool.params", string(paramsJSON)),
	)

	result, err := tool.Execute(ctx, params)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resultJSON, _ := json.Marshal(result)
	span.SetAttributes(
		attribute.String("gen_ai.tool.result", string(resultJSON)),
	)

	return result, nil
}

func RunToolMode() {
	fmt.Println("=== Tool调用模式示例 ===")

	toolService := NewToolService()
	ctx := context.Background()

	// 示例1: 天气查询
	fmt.Println("\n1. 查询天气:")
	weatherResult, err := toolService.ExecuteTool(ctx, "get_weather", map[string]interface{}{
		"city": "北京",
	})
	if err != nil {
		fmt.Printf("Weather tool error: %v\n", err)
	} else {
		resultJSON, _ := json.MarshalIndent(weatherResult, "", "  ")
		fmt.Printf("天气结果: %s\n", string(resultJSON))
	}

	// 示例2: 计算器
	fmt.Println("\n2. 数学计算:")
	calcResult, err := toolService.ExecuteTool(ctx, "calculator", map[string]interface{}{
		"operation": "multiply",
		"a":         15.5,
		"b":         2.0,
	})
	if err != nil {
		fmt.Printf("Calculator tool error: %v\n", err)
	} else {
		resultJSON, _ := json.MarshalIndent(calcResult, "", "  ")
		fmt.Printf("计算结果: %s\n", string(resultJSON))
	}
}
