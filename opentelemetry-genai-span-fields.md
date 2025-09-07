# GenAI Span 指标汇总

## 概述

本文档汇总了 OpenTelemetry 语义约定中与 GenAI 相关的 Span 指标，包括字段名、中英文解释、样例和枚举值。

## 通用字段

### 操作名称 (Operation Name)

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `gen_ai.operation.name` | string | 执行的操作名称 | The name of the operation being performed | `chat`, `generate_content`, `text_completion` | 见下方枚举值 |

**枚举值:**
- `chat` - 聊天完成操作，如 OpenAI Chat API
- `create_agent` - 创建 GenAI 代理
- `embeddings` - 嵌入操作，如 OpenAI Create embeddings API
- `execute_tool` - 执行工具
- `generate_content` - 多模态内容生成操作，如 Gemini Generate Content
- `invoke_agent` - 调用 GenAI 代理
- `text_completion` - 文本补全操作，如 OpenAI Completions API (Legacy)

### 提供商名称 (Provider Name)

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `gen_ai.provider.name` | string | 生成式 AI 提供商 | The Generative AI provider as identified by the client or server instrumentation | `openai`, `gcp.gen_ai`, `gcp.vertex_ai` | 见下方枚举值 |

**枚举值:**
- `anthropic` - Anthropic
- `aws.bedrock` - AWS Bedrock
- `azure.ai.inference` - Azure AI Inference
- `azure.ai.openai` - Azure OpenAI
- `cohere` - Cohere
- `deepseek` - DeepSeek
- `gcp.gemini` - Gemini (用于访问 'generativelanguage.googleapis.com' 端点)
- `gcp.gen_ai` - 任何 Google 生成式 AI 端点 (当具体后端未知时使用)
- `gcp.vertex_ai` - Vertex AI (用于访问 'aiplatform.googleapis.com' 端点)
- `groq` - Groq
- `ibm.watsonx.ai` - IBM Watsonx AI
- `mistral_ai` - Mistral AI
- `openai` - OpenAI
- `perplexity` - Perplexity
- `x_ai` - xAI

### 输出类型 (Output Type)

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `gen_ai.output.type` | string | 客户端请求的内容类型 | Represents the content type requested by the client | `text`, `json`, `image` | 见下方枚举值 |

**枚举值:**
- `image` - 图像
- `json` - 具有已知或未知模式的 JSON 对象
- `speech` - 语音
- `text` - 纯文本

## 请求相关字段

### 模型信息

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.request.model` | string | 请求的 GenAI 模型名称 | The name of the GenAI model a request is being made to | `gpt-4` |

### 请求参数

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.request.choice.count` | int | 要返回的候选完成目标数量 | The target number of candidate completions to return | `3` |
| `gen_ai.request.seed` | int | 相同种子值更可能返回相同结果 | Requests with same seed value more likely to return same result | `100` |
| `gen_ai.request.frequency_penalty` | double | GenAI 请求的频率惩罚设置 | The frequency penalty setting for the GenAI request | `0.1` |
| `gen_ai.request.max_tokens` | int | 模型为请求生成的最大令牌数 | The maximum number of tokens the model generates for a request | `100` |
| `gen_ai.request.presence_penalty` | double | GenAI 请求的存在惩罚设置 | The presence penalty setting for the GenAI request | `0.1` |
| `gen_ai.request.stop_sequences` | string[] | 模型将用来停止生成更多令牌的序列列表 | List of sequences that the model will use to stop generating further tokens | `["forest", "lived"]` |
| `gen_ai.request.temperature` | double | GenAI 请求的温度设置 | The temperature setting for the GenAI request | `0.0` |
| `gen_ai.request.top_k` | double | GenAI 请求的 top_k 采样设置 | The top_k sampling setting for the GenAI request | `1.0` |
| `gen_ai.request.top_p` | double | GenAI 请求的 top_p 采样设置 | The top_p sampling setting for the GenAI request | `1.0` |
| `gen_ai.request.encoding_formats` | string[] | 嵌入操作中请求的编码格式 | The encoding formats requested in an embeddings operation | `["base64"]` |

## 响应相关字段

### 响应标识

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.response.id` | string | 完成的唯一标识符 | The unique identifier for the completion | `chatcmpl-123` |
| `gen_ai.response.model` | string | 生成响应的模型名称 | The name of the model that generated the response | `gpt-4-0613` |
| `gen_ai.response.finish_reasons` | string[] | 模型停止生成令牌的原因数组 | Array of reasons the model stopped generating tokens | `["stop"]`, `["stop", "length"]` |

### 使用统计

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.usage.input_tokens` | int | GenAI 输入（提示）中使用的令牌数 | The number of tokens used in the GenAI input (prompt) | `100` |
| `gen_ai.usage.output_tokens` | int | GenAI 响应（完成）中使用的令牌数 | The number of tokens used in the GenAI response (completion) | `180` |

## 会话和消息

### 会话信息

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.conversation.id` | string | 对话（会话、线程）的唯一标识符 | The unique identifier for a conversation (session, thread) | `conv_5j66UpCpwteGg4YSxUnt7lPY` |

### 消息内容

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.input.messages` | any | 提供给模型的聊天历史记录 | The chat history provided to the model as an input | 见下方示例 |
| `gen_ai.output.messages` | any | 模型返回的消息 | Messages returned by the model | 见下方示例 |
| `gen_ai.system_instructions` | any | 单独提供给 GenAI 模型的系统消息或指令 | The system message or instructions provided to the GenAI model | 见下方示例 |

**消息格式示例:**
```json
// 输入消息示例
[
  {
    "role": "user",
    "parts": [
      {
        "type": "text",
        "content": "Weather in Paris?"
      }
    ]
  },
  {
    "role": "assistant",
    "parts": [
      {
        "type": "tool_call",
        "id": "call_VSPygqKTWdrhaFErNvMV18Yl",
        "name": "get_weather",
        "arguments": {
          "location": "Paris"
        }
      }
    ]
  }
]

// 输出消息示例
[
  {
    "role": "assistant",
    "parts": [
      {
        "type": "text",
        "content": "The weather in Paris is currently rainy with a temperature of 57°F."
      }
    ],
    "finish_reason": "stop"
  }
]
```

## 工具相关字段

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `gen_ai.tool.call.id` | string | 工具调用标识符 | The tool call identifier | `call_mszuSIzqtI65i1wAUOE8w5H4` | - |
| `gen_ai.tool.description` | string | 工具描述 | The tool description | `Multiply two numbers` | - |
| `gen_ai.tool.name` | string | 代理使用的工具名称 | Name of the tool utilized by the agent | `Flights` | - |
| `gen_ai.tool.type` | string | 代理使用的工具类型 | Type of the tool utilized by the agent | `function`, `extension`, `datastore` | 见下方枚举值 |

**工具类型枚举值:**
- `function` - 在客户端执行的函数，代理生成参数，客户端执行逻辑
- `extension` - 在代理端执行的工具，直接调用外部 API
- `datastore` - 代理用于访问和查询结构化或非结构化外部数据的工具

## 代理相关字段

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `gen_ai.agent.description` | string | 应用提供的 GenAI 代理的描述 | Free-form description of the GenAI agent provided by the application | `Helps with math problems` |
| `gen_ai.agent.id` | string | GenAI 代理的唯一标识符 | The unique identifier of the GenAI agent | `asst_5j66UpCpwteGg4YSxUnt7lPY` |
| `gen_ai.agent.name` | string | 应用提供的 GenAI 代理的可读名称 | Human-readable name of the GenAI agent provided by the application | `Math Tutor` |
| `gen_ai.data_source.id` | string | 数据源标识符 | The data source identifier | `H7STPQYOND` |

## 服务器和错误信息

### 服务器信息

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `server.address` | string | GenAI 服务器地址 | GenAI server address | `example.com`, `10.1.2.80` |
| `server.port` | int | GenAI 服务器端口 | GenAI server port | `80`, `8080`, `443` |

### 错误信息

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `error.type` | string | 操作结束的错误类别描述 | Describes a class of error the operation ended with | `timeout`, `500` | 见下方枚举值 |

**错误类型枚举值:**
- `_OTHER` - 当仪器未定义自定义值时使用的回退错误值

## 特定提供商字段

### AWS Bedrock 特有字段

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `aws.bedrock.guardrail.id` | string | AWS Bedrock Guardrail 的唯一标识符 | The unique identifier of the AWS Bedrock Guardrail | `sgi5gkybzqak` |
| `aws.bedrock.knowledge_base.id` | string | AWS Bedrock 知识库的唯一标识符 | The unique identifier of the AWS Bedrock Knowledge base | `XFWUPB9PAW` |

### Azure AI Inference 特有字段

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 |
|--------|------|----------|----------|------|
| `azure.resource_provider.namespace` | string | 客户端识别的 Azure 资源提供程序命名空间 | Azure Resource Provider Namespace as recognized by the client | `Microsoft.CognitiveServices` |

### OpenAI 特有字段

| 字段名 | 类型 | 中文解释 | 英文解释 | 样例 | 枚举值 |
|--------|------|----------|----------|------|--------|
| `openai.request.service_tier` | string | 请求的服务层级 | The service tier requested | `auto`, `default` | 见下方枚举值 |
| `openai.response.service_tier` | string | 响应使用的服务层级 | The service tier used for the response | `scale`, `default` | 见下方枚举值 |
| `openai.response.system_fingerprint` | string | 跟踪生成式 AI 环境任何变化的指纹 | A fingerprint to track any eventual change in the Generative AI environment | `fp_44709d6fcb` | - |

**服务层级枚举值:**
- `auto` - 系统将使用扩展层级积分，直到用完为止
- `default` - 系统将使用默认扩展层级

## 指标汇总

### 客户端指标

#### 令牌使用量 (Token Usage)

| 指标名 | 类型 | 单位 | 描述 |
|--------|------|------|------|
| `gen_ai.client.token.usage` | Histogram | `{token}` | 使用的输入和输出令牌数量 |

**必需属性:**
- `gen_ai.operation.name`
- `gen_ai.provider.name`
- `gen_ai.token.type`

**条件必需属性:**
- `gen_ai.request.model` (如果可用)

**推荐属性:**
- `gen_ai.response.model`
- `server.address`

**令牌类型枚举值:**
- `input` - 输入令牌（提示、输入等）
- `output` - 输出令牌（完成、响应等）

#### 操作持续时间 (Operation Duration)

| 指标名 | 类型 | 单位 | 描述 |
|--------|------|------|------|
| `gen_ai.client.operation.duration` | Histogram | `s` | GenAI 操作持续时间 |

**必需属性:**
- `gen_ai.operation.name`
- `gen_ai.provider.name`

**条件必需属性:**
- `gen_ai.request.model` (如果可用)
- `error.type` (如果操作以错误结束)

### 服务器指标

#### 请求持续时间 (Request Duration)

| 指标名 | 类型 | 单位 | 描述 |
|--------|------|------|------|
| `gen_ai.server.request.duration` | Histogram | `s` | 生成式 AI 服务器请求持续时间 |

**必需属性:**
- `gen_ai.operation.name`
- `gen_ai.provider.name`

#### 每输出令牌时间 (Time Per Output Token)

| 指标名 | 类型 | 单位 | 描述 |
|--------|------|------|------|
| `gen_ai.server.time_per_output_token` | Histogram | `s` | 成功响应后每个输出令牌的生成时间 |

**必需属性:**
- `gen_ai.operation.name`
- `gen_ai.provider.name`

#### 首令牌时间 (Time To First Token)

| 指标名 | 类型 | 单位 | 描述 |
|--------|------|------|------|
| `gen_ai.server.time_to_first_token` | Histogram | `s` | 生成响应第一个令牌的时间 |

**必需属性:**
- `gen_ai.operation.name`
- `gen_ai.provider.name`

## 事件汇总

### 推理操作详情事件

| 事件名 | 描述 | 必需属性 | 推荐属性 | 可选属性 |
|--------|------|----------|----------|----------|
| `event.gen_ai.client.inference.operation.details` | 描述 GenAI 完成请求的详情 | `gen_ai.operation.name` | 见下方属性列表 | 见下方属性列表 |

**事件属性:**
- 必需: `gen_ai.operation.name`
- 条件必需: `error.type` (如果操作以错误结束), `gen_ai.conversation.id` (如果可用), `gen_ai.output.type` (如果适用)
- 推荐: `gen_ai.request.choice.count`, `gen_ai.request.model`, `gen_ai.request.seed`, `server.port`, `gen_ai.request.frequency_penalty`, `gen_ai.request.max_tokens`, `gen_ai.request.presence_penalty`, `gen_ai.request.stop_sequences`, `gen_ai.request.temperature`, `gen_ai.request.top_p`, `gen_ai.response.finish_reasons`, `gen_ai.response.id`, `gen_ai.response.model`, `gen_ai.usage.input_tokens`, `gen_ai.usage.output_tokens`, `server.address`
- 可选: `gen_ai.input.messages`, `gen_ai.output.messages`, `gen_ai.system_instructions`
