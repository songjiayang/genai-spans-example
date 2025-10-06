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
		fmt.Println("  go run main.go chat                    # 运行聊天模式示例 (console导出器)")
		fmt.Println("  go run main.go tool                    # 运行工具调用模式示例 (console导出器)")
		fmt.Println("  go run main.go agent                   # 运行Agent模式示例 (console导出器)")
		fmt.Println("  go run main.go chat --http             # 运行聊天模式示例 (HTTP导出器)")
		fmt.Println("  go run main.go tool --http             # 运行工具调用模式示例 (HTTP导出器)")
		fmt.Println("  go run main.go agent --http            # 运行Agent模式示例 (HTTP导出器)")
		fmt.Println("")
		fmt.Println("环境变量:")
		fmt.Println("  OTEL_EXPORTER_OTLP_ENDPOINT            # OTLP端点 (默认: http://localhost:4318)")
		fmt.Println("  OTEL_SERVICE_NAME                      # 服务名称 (默认: gen-ai-example)")
		fmt.Println("  OTEL_TRACES_EXPORTER                   # 导出器类型 (console/http/otlp/auto)")
		return
	}

	// 检查是否使用HTTP导出器
	useHTTP := false
	mode := os.Args[1]

	for i, arg := range os.Args {
		if arg == "--http" {
			useHTTP = true
			// 移除--http参数
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	// 初始化telemetry
	var cleanup func()

	if useHTTP {
		// 如果指定了--http，强制使用HTTP导出器
		config := telemetry.Config{
			ExporterType: telemetry.ExporterHTTP,
			Endpoint:     "http://localhost:4318",
			ServiceName:  "gen-ai-example",
		}
		cleanup = telemetry.InitTracerWithConfig(config)
		fmt.Println("使用HTTP导出器，端点: http://localhost:4318")
	} else {
		// 否则从环境变量获取配置
		config := telemetry.GetConfigFromEnv()
		cleanup = telemetry.InitTracerWithConfig(config)

		if config.ExporterType == telemetry.ExporterHTTP {
			fmt.Printf("使用HTTP导出器，端点: %s\n", config.Endpoint)
		} else {
			fmt.Println("使用console导出器")
		}
	}
	defer cleanup()

	// 运行相应的模式
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
