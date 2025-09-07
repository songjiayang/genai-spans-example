package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"gen-ai-example/telemetry"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (containsIgnoreCase(s, substr) ||
			containsWord(s, substr)))
}

// containsIgnoreCase 不区分大小写包含检查
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// containsWord 检查是否包含完整单词
func containsWord(s, word string) bool {
	words := strings.Fields(s)
	for _, w := range words {
		if strings.Contains(strings.ToLower(w), strings.ToLower(word)) {
			return true
		}
	}
	return false
}

type ChatRequest struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type ChatResponse struct {
	Reply     string    `json:"reply"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatService struct {
	tracer trace.Tracer
}

func NewChatService() *ChatService {
	return &ChatService{
		tracer: telemetry.GetTracer("chat-service"),
	}
}

func (cs *ChatService) ProcessChat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	conversationID := uuid.New().String()
	_, span := cs.tracer.Start(ctx, "chat.process",
		trace.WithAttributes(
			semconv.GenAIOperationNameChat,
			semconv.GenAIProviderNameOpenAI,
			semconv.GenAIRequestModel("gpt-3.5-turbo"),
			semconv.GenAIConversationID(conversationID),
			semconv.GenAIInputMessagesKey.String(fmt.Sprintf(`[{"role":"user","content":"%s"}]`, req.Message)),
		),
	)
	defer span.End()

	time.Sleep(100 * time.Millisecond)

	// 根据用户消息生成合适的AI回复
	var reply string
	message := req.Message

	// 检查多个关键词，优先级高的先匹配
	switch {
	case contains(message, "Go语言") || contains(message, "Golang") || (contains(message, "Go") && contains(message, "语言")):
		reply = "Go语言是Google开发的一种静态强类型、编译型语言。它具有简洁的语法、高效的并发处理能力和优秀的性能，非常适合构建网络服务和分布式系统。"
	case contains(message, "天气"):
		reply = "我无法获取实时天气信息，但您可以使用天气查询工具来获取准确的天气数据。"
	case contains(message, "谢谢"):
		reply = "不客气！如果您还有其他问题，随时告诉我。"
	case contains(message, "你好") || contains(message, "您好"):
		reply = "你好！很高兴为您服务。我是一个AI助手，可以回答您的问题和提供帮助。"
	default:
		reply = fmt.Sprintf("我理解您说的是：%s。这是一个很有趣的话题，我可以为您提供更多相关信息。", message)
	}

	response := &ChatResponse{
		Reply:     reply,
		Timestamp: time.Now(),
	}

	span.SetAttributes(
		semconv.GenAIOutputMessagesKey.String(fmt.Sprintf(`[{"role":"assistant","content":"%s"}]`, response.Reply)),
		semconv.GenAIUsageOutputTokens(len(response.Reply)),
		semconv.GenAIUsageInputTokens(len(req.Message)),
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

	return response, nil
}

func RunChatMode() {
	fmt.Println("=== 通用AI Chat模式示例 ===")

	chatService := NewChatService()
	ctx := context.Background()
	req := ChatRequest{
		Message: "你好，请介绍一下Go语言",
		UserID:  "user123",
	}

	response, err := chatService.ProcessChat(ctx, req)
	if err != nil {
		fmt.Printf("Chat processing failed: %v\n", err)
		return
	}

	fmt.Printf("用户消息: %s\n", req.Message)
	fmt.Printf("AI回复: %s\n", response.Reply)
	fmt.Printf("时间戳: %s\n", response.Timestamp.Format(time.RFC3339))
}
