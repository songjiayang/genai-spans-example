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

	time.Sleep(10 * time.Millisecond)

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

// ChatModelResponse 模拟聊天模型的响应
type ChatModelResponse struct {
	Role    string                   `json:"role"`
	Content string                   `json:"content"`
	Tools   []map[string]interface{} `json:"tools"`
}

// SimulateChatModelCall 模拟调用聊天模型来获取工具调用决策
func (ts *ToolService) SimulateChatModelCall(ctx context.Context, userMessage string) (*ChatModelResponse, error) {
	// 创建聊天模型调用追踪
	conversationID := uuid.New().String()
	_, span := ts.tracer.Start(ctx, "chat-model.call",
		trace.WithAttributes(
			semconv.GenAIOperationNameChat,
			semconv.GenAIProviderNameOpenAI,
			semconv.GenAIRequestModel("gpt-3.5-turbo"),
			semconv.GenAIConversationID(conversationID),
			semconv.GenAIInputMessagesKey.String(fmt.Sprintf(`[{"role":"user","content":"%s"}]`, userMessage)),
		),
	)
	defer span.End()

	// 模拟模型思考时间
	time.Sleep(30 * time.Millisecond)

	// 模拟工具调用决策
	tools := []map[string]interface{}{
		{
			"name":        "get_weather",
			"description": "Get weather information for a specified city",
		},
		{
			"name":        "calculator",
			"description": "Perform basic mathematical calculations",
		},
	}

	// 根据用户消息决定需要哪些工具
	if userMessage == "查询北京的天气，然后计算10+25的结果" {
		tools = []map[string]interface{}{
			{
				"name":        "get_weather",
				"description": "Get weather information for a specified city",
			},
			{
				"name":        "calculator",
				"description": "Perform basic mathematical calculations",
			},
		}
	}

	resp := ChatModelResponse{
		Role:    "assistant",
		Content: "我需要调用一些工具来帮助您完成请求",
		Tools:   tools,
	}

	respJson, _ := json.Marshal([]ChatModelResponse{resp})

	span.SetAttributes(
		semconv.GenAIOutputMessagesKey.String(string(respJson)),
		semconv.GenAIUsageOutputTokens(len(string(respJson))),
		semconv.GenAIUsageInputTokens(len(userMessage)),
		semconv.GenAIResponseID(fmt.Sprintf("chatcmpl-%d", time.Now().Unix())),
		semconv.GenAIResponseFinishReasons("stop"),
		semconv.GenAIRequestMaxTokens(2048),
		semconv.GenAIRequestTemperature(0.7),
		semconv.GenAIRequestTopP(1.0),
		semconv.GenAIRequestFrequencyPenalty(0),
		semconv.GenAIRequestPresencePenalty(0),
		semconv.GenAIRequestChoiceCount(1),
		semconv.GenAIRequestSeed(42),
		semconv.GenAIOutputTypeText,
	)

	return &resp, nil
}

// ExecuteToolChain 执行一系列工具调用，共享同一个 trace ID
func (ts *ToolService) ExecuteToolChain(ctx context.Context, userMessage string) (map[string]interface{}, error) {
	// 在根span的上下文中模拟调用聊天模型
	chatResponse, err := ts.SimulateChatModelCall(ctx, userMessage)
	if err != nil {
		return nil, err
	}

	results := make(map[string]interface{})

	// 记录聊天模型响应
	chatJSON, _ := json.Marshal(chatResponse)
	fmt.Printf("🤖 模型响应: %s\n", string(chatJSON))

	// 执行每个工具调用
	for _, tool := range chatResponse.Tools {
		toolName := tool["name"].(string)

		switch toolName {
		case "get_weather":
			fmt.Printf("🌤️  正在查询天气...\n")
			weatherResult, err := ts.ExecuteTool(ctx, "get_weather", map[string]interface{}{
				"city": "北京",
			})

			if err != nil {
				fmt.Printf("❌ 天气查询失败: %v\n", err)
			} else {
				results["weather"] = weatherResult
				weatherJSON, _ := json.MarshalIndent(weatherResult, "", "  ")
				fmt.Printf("✅ 天气结果: %s\n", string(weatherJSON))
			}

		case "calculator":
			fmt.Printf("🧮 正在执行计算...\n")
			calcResult, err := ts.ExecuteTool(ctx, "calculator", map[string]interface{}{
				"operation": "add",
				"a":         10.0,
				"b":         25.0,
			})

			if err != nil {
				fmt.Printf("❌ 计算失败: %v\n", err)
			} else {
				results["calculator"] = calcResult
				calcJSON, _ := json.MarshalIndent(calcResult, "", "  ")
				fmt.Printf("✅ 计算结果: %s\n", string(calcJSON))
			}
		}
	}

	return results, nil
}

func (ts *ToolService) RegisterTool(tool Tool) {
	ts.tools[tool.Name()] = tool
}

func (ts *ToolService) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	ctx, span := ts.tracer.Start(ctx, "tool.execute",
		trace.WithAttributes(
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

	// 创建根上下文和追踪
	ctx, rootSpan := telemetry.GetTracer("tool-mode").Start(context.Background(), "tool-mode.root")
	defer rootSpan.End()

	toolService := NewToolService()

	// 示例1: 模拟完整的工具链调用（模拟聊天模型 + 工具调用）
	fmt.Println("\n1. 模拟聊天模型 + 工具调用:")
	fmt.Println("用户消息: 查询北京的天气，然后计算10+25的结果")

	results, err := toolService.ExecuteToolChain(ctx, "查询北京的天气，然后计算10+25的结果")
	if err != nil {
		fmt.Printf("工具链调用失败: %v\n", err)
	} else {
		fmt.Println("\n📊 综合结果:")
		resultsJSON, _ := json.MarshalIndent(results, "", "  ")
		fmt.Printf("%s\n", string(resultsJSON))
	}
}
