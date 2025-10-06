# Gen AI 示例项目

一个综合性示例项目，展示AI代理和工具集成与OpenTelemetry遥测功能。

## 功能特性

- **代理模式**: AI代理进行多步骤任务规划和执行
- **工具模式**: 单个工具执行，包含天气查询和计算器功能
- **遥测监控**: OpenTelemetry集成用于追踪和监控
- **结构化输出**: JSON格式结果和全面日志记录
- **多导出器支持**: Console和HTTP导出器，支持生产环境部署

## 项目结构

```
pkg/
├── agent/          # 代理实现，用于任务规划和执行
├── tool/           # 工具定义和执行服务
└── telemetry/      # OpenTelemetry配置
```

## 组件说明

### 代理包 (`pkg/agent/`)

代理系统提供智能任务规划和执行功能：

- **任务规划**: 自动将用户目标分解为可执行任务
- **任务执行**: 按顺序运行任务，具备适当的错误处理
- **工具集成**: 无缝调用注册的工具来实现目标
- **遥测追踪**: 代理操作的全面追踪

**关键函数:**
- `PlanTasks()`: 分析用户目标并创建任务计划
- `ExecuteTasks()`: 执行计划任务，具备适当的错误处理
- `RegisterTool()`: 注册可用工具供代理使用

### 工具包 (`pkg/tool/`)

工具系统提供可扩展的工具执行功能：

- **天气工具**: 获取指定城市的天气信息
- **计算器工具**: 执行基本数学运算
- **工具服务**: 集中式工具注册和执行
- **遥测追踪**: 工具调用的详细追踪

**可用工具:**
- `get_weather`: "获取指定城市的天气信息"
- `calculator`: "执行基本数学计算"

### 遥测包 (`pkg/telemetry/`)

OpenTelemetry集成提供全面的可观测性：

- **追踪**: 所有操作的分布式追踪
- **多导出器支持**: Console和HTTP导出器，支持生产和开发环境
- **自动检测**: 根据环境变量自动选择最佳导出器
- **语义约定**: 使用OpenAI GenAI语义约定
- **自定义属性**: AI操作增强的元数据

## 使用方法

### 基本命令

项目支持多种运行模式和遥测配置：

```bash
# 聊天模式示例 (默认console导出器)
go run main.go chat

# 工具调用模式示例 (默认console导出器)
go run main.go tool

# 代理模式示例 (默认console导出器)
go run main.go agent

# 强制使用HTTP导出器
go run main.go chat --http
go run main.go tool --http
go run main.go agent --http
```

### 代理模式

运行代理示例查看多步骤任务规划和执行：

```bash
go run main.go agent
```

示例输出：
```
=== Agent模式示例 ===

目标: 请帮我查询北京的天气，然后计算10+25的结果

1. 任务规划阶段:
  - task-1: 查询天气信息
  - task-2: 执行数学计算
  - task-final: 总结执行结果

2. 任务执行阶段:
3. 执行结果:
```

### 工具模式

运行单个工具示例：

```bash
go run main.go tool
```

示例输出：
```
=== Tool调用模式示例 ===

1. 查询天气:
天气结果: {
  "city": "北京",
  "temperature": "22°C",
  "condition": "晴天",
  "humidity": "65%"
}

2. 数学计算:
计算结果: {
  "operation": "multiply",
  "a": 15.5,
  "b": 2.0,
  "result": 31
}
```

## 遥测功能特性

项目包含AI操作的全面遥测功能：

### OpenTelemetry GenAI Span 字段参考

详细字段说明和规范请参考：[OpenTelemetry GenAI Span 字段文档](opentelemetry-genai-span-fields.md)

### 追踪属性

**代理操作:**
- `gen_ai.agent.planned_tasks_count`: 计划任务数量
- `gen_ai.input.tokens`: 输入令牌数量
- `gen_ai.output.tokens`: 输出令牌数量
- `gen_ai.agent.name`: 代理名称
- `gen_ai.agent.description`: 代理描述

**工具操作:**
- `gen_ai.tool.name`: 工具名称
- `gen_ai.tool.description`: 工具描述
- `gen_ai.tool.type`: 工具类型
- `gen_ai.tool.params`: 工具参数
- `gen_ai.tool.result`: 工具执行结果

**消息追踪:**
- `gen_ai.input.messages`: JSON格式的输入消息
- `gen_ai.output.messages`: JSON格式的输出消息

### 令牌计数

自动令牌估算：
- **输入令牌**: 基于用户目标长度
- **输出令牌**: 基于计划任务复杂度

### 语义约定

使用OpenAI GenAI语义约定：
- `gen_ai.provider.name`: "openai"
- `gen_ai.operation.name`: 操作描述
- `gen_ai.tool.call_id`: 唯一调用标识符
- `gen_ai.usage.input_tokens`: 输入令牌数量
- `gen_ai.usage.output_tokens`: 输出令牌数量

## 开发指南

### 添加新工具

1. 实现 `Tool` 接口：
```go
type NewTool struct{}

func (t *NewTool) Name() string {
    return "new_tool"
}

func (t *NewTool) Description() string {
    return "新工具的描述"
}

func (t *NewTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // 工具实现
}
```

2. 向代理注册：
```go
agent.RegisterTool(&NewTool{})
```

3. 在 `PlanTasks()` 中添加任务规划逻辑

### 扩展遥测功能

向span添加自定义属性：
```go
span.SetAttributes(
    attribute.String("custom.attribute", "值"),
    attribute.Int("custom.metric", 42),
)
```

## 配置

### 遥测配置

项目支持多种OpenTelemetry导出器配置方式：

#### 命令行选项
- `--http`: 强制使用HTTP导出器，连接到 localhost:4318

#### 环境变量
```bash
# 设置导出器类型 (console/http/otlp/auto)
export OTEL_TRACES_EXPORTER=http

# 设置OTLP端点
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

# 设置服务名称
export OTEL_SERVICE_NAME=gen-ai-example
```

#### 导出器类型说明
- **console**: 开发环境，输出到控制台
- **http**: 生产环境，发送到OTLP兼容的收集器
- **auto**: 自动检测，优先HTTP，失败时回退到console

#### 依赖项

主要依赖项：
- `go.opentelemetry.io/otel`: OpenTelemetry API
- `go.opentelemetry.io/otel/exporters/stdout/stdouttrace`: Console导出器
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`: HTTP导出器
- `go.opentelemetry.io/otel/semconv/v1.24.0`: 语义约定
- `github.com/google/uuid`: UUID生成

## 测试

运行项目查看示例演示：

### 开发环境测试（Console导出器）
```bash
go run main.go agent
go run main.go tool
go run main.go chat
```

### 生产环境测试（HTTP导出器）
```bash
# 强制使用HTTP导出器
go run main.go agent --http

# 使用环境变量配置
export OTEL_TRACES_EXPORTER=http
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
go run main.go agent
```

示例演示：
1. 代理任务规划，包含天气和计算器任务
2. 工具执行，具备适当的错误处理
3. 全面遥测追踪（输出到console或HTTP）
4. 结构化JSON输出格式
5. 多导出器自动检测和回退机制

## 许可证

本项目用于演示目的，包含AI代理和工具集成模式的教育示例。