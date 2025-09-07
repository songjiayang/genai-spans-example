package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gen-ai-example/pkg/tool"
	"gen-ai-example/telemetry"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Task struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Params      map[string]interface{} `json:"params"`
	Status      string                 `json:"status"`
	Result      interface{}            `json:"result,omitempty"`
}

type Agent struct {
	name   string
	tracer trace.Tracer
	tasks  []Task
	tools  map[string]tool.Tool
}

func NewAgent(name string) *Agent {
	return &Agent{
		name:   name,
		tracer: telemetry.GetTracer(fmt.Sprintf("agent-%s", name)),
		tasks:  make([]Task, 0),
		tools:  make(map[string]tool.Tool),
	}
}

func (a *Agent) RegisterTool(tool tool.Tool) {
	a.tools[tool.Name()] = tool
}

// generatePlanMessages 生成任务规划结果的AI输出消息
func generatePlanMessages(tasks []Task, objective string) string {
	var taskDescriptions []string
	for _, task := range tasks {
		taskDescriptions = append(taskDescriptions, fmt.Sprintf(`{"task_id":"%s","description":"%s","type":"%s"}`, task.ID, task.Description, task.Type))
	}

	tasksJSON := fmt.Sprintf("[%s]", strings.Join(taskDescriptions, ","))

	return fmt.Sprintf(`[{"role":"assistant","content":"Task planning completed for objective: %s","tasks":%s}]`, objective, tasksJSON)
}

// estimateInputTokens 估算输入tokens数量
func estimateInputTokens(objective string) int {
	// 简单估算：每个中文字符算1个token，每个英文字符算0.25个token
	tokens := 0
	for _, r := range objective {
		if r > 127 {
			tokens += 1 // 中文等Unicode字符
		} else {
			tokens += 1 // 英文字符
		}
	}
	return tokens + 10 // 添加消息结构的固定开销
}

// estimateOutputTokens 估算输出tokens数量
func estimateOutputTokens(tasks []Task) int {
	// 简单估算：每个任务大约20个tokens
	tokens := len(tasks) * 20
	return tokens + 30 // 添加响应结构的固定开销
}

func (a *Agent) PlanTasks(ctx context.Context, objective string) error {
	_, span := a.tracer.Start(ctx, "agent.plan_tasks",
		trace.WithAttributes(
			semconv.GenAIProviderNameOpenAI,
			semconv.GenAIOperationNameCreateAgent,
			semconv.GenAIAgentID(uuid.NewString()),
			semconv.GenAIAgentName(a.name),
			semconv.GenAIAgentDescription(objective),
		),
	)
	defer span.End()

	time.Sleep(200 * time.Millisecond)

	var tasks []Task
	if strings.Contains(strings.ToLower(objective), "天气") {
		tasks = append(tasks, Task{
			ID:          "task-1",
			Description: "查询天气信息",
			Type:        "tool_call",
			Params: map[string]interface{}{
				"tool": "get_weather",
				"city": "北京",
			},
			Status: "pending",
		})
	}

	if strings.Contains(strings.ToLower(objective), "计算") {
		tasks = append(tasks, Task{
			ID:          "task-2",
			Description: "执行数学计算",
			Type:        "tool_call",
			Params: map[string]interface{}{
				"tool":      "calculator",
				"operation": "add",
				"a":         10.0,
				"b":         25.0,
			},
			Status: "pending",
		})
	}

	tasks = append(tasks, Task{
		ID:          "task-final",
		Description: "总结执行结果",
		Type:        "summarize",
		Params:      map[string]interface{}{},
		Status:      "pending",
	})

	a.tasks = tasks

	// 生成任务规划结果的AI输出消息
	outputMessages := generatePlanMessages(tasks, objective)

	// 估算输入输出tokens
	inputTokens := estimateInputTokens(objective)
	outputTokens := estimateOutputTokens(tasks)

	span.SetAttributes(
		attribute.Int("gen_ai.agent.planned_tasks_count", len(tasks)),
		semconv.GenAIInputMessagesKey.String(fmt.Sprintf(`[{"role":"user","content":"%s"}]`, objective)),
		semconv.GenAIOutputMessagesKey.String(outputMessages),
		semconv.GenAIUsageInputTokens(inputTokens),
		semconv.GenAIUsageOutputTokens(outputTokens),
	)

	return nil
}

func (a *Agent) ExecuteTasks(ctx context.Context) error {
	ctx, span := a.tracer.Start(ctx, "agent.execute_tasks",
		trace.WithAttributes(
			semconv.GenAIProviderNameOpenAI,
			semconv.GenAIOperationNameInvokeAgent,
			semconv.GenAIAgentName(a.name),
			attribute.Int("gen_id.agent.total_tasks", len(a.tasks)),
		),
	)
	defer span.End()

	var results []interface{}

	for i := range a.tasks {
		task := &a.tasks[i]
		if task.Status != "pending" {
			continue
		}

		taskCtx, taskSpan := a.tracer.Start(ctx, fmt.Sprintf("agent.execute_task.%s", task.ID),
			trace.WithAttributes(
				attribute.String("gen_id.agent.task_id", task.ID),
				attribute.String("gen_id.agent.task_type", task.Type),
				attribute.String("gen_id.agent.description", task.Description),
				semconv.GenAIToolName(fmt.Sprintf("task_%s", task.ID)),
				semconv.GenAIToolDescription(task.Description),
			),
		)

		var err error
		switch task.Type {
		case "tool_call":
			task.Result, err = a.executeToolTask(taskCtx, task)
		case "summarize":
			task.Result, err = a.executeSummaryTask(taskCtx, results)
		default:
			err = fmt.Errorf("unknown task type: %s", task.Type)
		}

		if err != nil {
			task.Status = "failed"
			taskSpan.RecordError(err)
		} else {
			task.Status = "completed"
			results = append(results, task.Result)
			resultJSON, _ := json.Marshal(task.Result)
			taskSpan.SetAttributes(
				attribute.String("gen_ai.task.result", string(resultJSON)),
			)
		}

		taskSpan.End()

		if err != nil {
			return fmt.Errorf("task %s failed: %w", task.ID, err)
		}
	}

	span.SetAttributes(
		attribute.Int("agent.completed_tasks", len(results)),
		semconv.GenAIOutputMessagesKey.String(fmt.Sprintf(`[{"role":"assistant","content":"任务执行完成，共完成%d个任务"}]`, len(results))),
	)

	return nil
}

func (a *Agent) executeToolTask(ctx context.Context, task *Task) (interface{}, error) {
	toolName, ok := task.Params["tool"].(string)
	if !ok {
		return nil, fmt.Errorf("missing tool name")
	}

	tool, exists := a.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	toolParams := make(map[string]interface{})
	for k, v := range task.Params {
		if k != "tool" {
			toolParams[k] = v
		}
	}

	return tool.Execute(ctx, toolParams)
}

func (a *Agent) executeSummaryTask(_ context.Context, results []interface{}) (interface{}, error) {
	time.Sleep(100 * time.Millisecond)

	summary := map[string]interface{}{
		"summary":      "任务执行完成",
		"total_tasks":  len(results),
		"results":      results,
		"completed_at": time.Now().Format(time.RFC3339),
	}

	return summary, nil
}

func (a *Agent) GetTaskResults() []Task {
	return a.tasks
}

func RunAgentMode() {
	fmt.Println("=== Agent模式示例 ===")

	agent := NewAgent("assistant")
	agent.RegisterTool(&tool.WeatherTool{})
	agent.RegisterTool(&tool.CalculatorTool{})

	ctx := context.Background()

	// 设定目标
	objective := "请帮我查询北京的天气，然后计算10+25的结果"
	fmt.Printf("目标: %s\n\n", objective)

	// 任务规划
	fmt.Println("1. 任务规划阶段:")
	err := agent.PlanTasks(ctx, objective)
	if err != nil {
		fmt.Printf("Task planning failed: %v\n", err)
		return
	}

	tasks := agent.GetTaskResults()
	for _, task := range tasks {
		fmt.Printf("  - %s: %s\n", task.ID, task.Description)
	}

	// 任务执行
	fmt.Println("\n2. 任务执行阶段:")
	err = agent.ExecuteTasks(ctx)
	if err != nil {
		fmt.Printf("Task execution failed: %v\n", err)
		return
	}

	// 显示结果
	fmt.Println("\n3. 执行结果:")
	tasks = agent.GetTaskResults()
	for _, task := range tasks {
		fmt.Printf("\n任务 %s (%s):\n", task.ID, task.Status)
		if task.Result != nil {
			resultJSON, _ := json.MarshalIndent(task.Result, "  ", "  ")
			fmt.Printf("  结果: %s\n", string(resultJSON))
		}
	}
}
