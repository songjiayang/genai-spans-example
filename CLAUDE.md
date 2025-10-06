# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based AI Agent and Tool Integration Example with OpenTelemetry observability. It demonstrates multi-step task planning, tool execution, and comprehensive telemetry following GenAI semantic conventions.

## Development Commands

### Build and Run
```bash
# Chat mode with rule-based conversation
go run main.go chat

# Tool execution mode (weather, calculator)
go run main.go tool

# Agent planning and execution mode
go run main.go agent

# Direct component execution
go run pkg/tool/tool.go
```

### Testing
No test framework currently configured. Use standard Go testing:
```bash
go test ./...
```

## Architecture Overview

### Core Components

**Agent System** (`pkg/agent/`):
- Multi-step task planning and execution
- Tool registry and orchestration
- OpenTelemetry tracing with GenAI semantic conventions
- Token usage estimation and cost tracking

**Tool System** (`pkg/tool/`):
- WeatherTool: Mock weather data service
- CalculatorTool: Mathematical operations
- ToolService: Centralized execution and registry
- Full telemetry instrumentation

**Chat System** (`pkg/chat/`):
- Rule-based keyword matching responses
- Multi-language support (Chinese/English)
- OpenTelemetry instrumentation

**Telemetry Layer** (`telemetry/`):
- OpenTelemetry tracer initialization
- Console exporter for development
- GenAI semantic convention compliance

### Design Patterns

- **Interface-driven** design with extensible Tool interface
- **Service pattern** for centralized business logic
- **Observer pattern** via OpenTelemetry tracing
- **Strategy pattern** for different execution modes

## Development Guidelines

### Code Structure
- Entry point: `main.go` with command-line routing
- Modular packages under `pkg/` for clean separation
- Comprehensive telemetry in all components
- Structured JSON outputs for all operations

### Telemetry Standards
Follow OpenTelemetry GenAI semantic conventions:
- Use `genai.system`, `genai.request.model`, `genai.usage.input_tokens`, etc.
- Span names should describe the operation clearly
- Include error attributes when spans fail

### Tool Development
- Implement `Tool` interface with `Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)`
- Add telemetry spans with appropriate attributes
- Return structured JSON responses
- Handle errors gracefully with context cancellation

### Agent Planning
- Use multi-step task decomposition
- Register tools with clear descriptions
- Implement token estimation for cost tracking
- Handle tool failures and retry logic

## Key Files and Entry Points

- `main.go` - Command-line router and telemetry initialization
- `pkg/agent/agent.go` - Core agent implementation
- `pkg/tool/tool.go` - Tool registry and execution
- `pkg/chat/chat.go` - Chat service implementation
- `telemetry/telemetry.go` - OpenTelemetry configuration

## Documentation References

- `README.md` - Comprehensive Chinese documentation
- `opentelemetry-genai-span-fields.md` - Detailed OpenTelemetry GenAI semantic convention reference

## Important Notes

- This project primarily uses Chinese for documentation and user-facing content
- All components include comprehensive OpenTelemetry tracing
- Tool execution supports both direct calls and agent orchestration
- The codebase follows Go best practices with proper error handling and context usage