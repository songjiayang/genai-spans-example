package main

import (
	"fmt"
	"os"
	"time"

	"gen-ai-example/pkg/agent"
	"gen-ai-example/pkg/chat"
	"gen-ai-example/pkg/tool"
	"gen-ai-example/telemetry"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法:")
		fmt.Println("  go run main.go chat    # 运行聊天模式示例")
		fmt.Println("  go run main.go tool    # 运行工具调用模式示例")
		fmt.Println("  go run main.go agent   # 运行Agent模式示例")
		return
	}

	// 初始化telemetry
	cleanup := telemetry.InitTracer()
	defer cleanup()

	mode := os.Args[1]
	switch mode {
	case "chat":
		chat.RunChatMode()
	case "tool":
		tool.RunToolMode()
	case "agent":
		agent.RunAgentMode()
	default:
		fmt.Printf("未知模式: %s\n", mode)
		return
	}

	// 等待trace输出
	time.Sleep(1 * time.Second)
}
