# Grafana Tempo 设置指南

本指南介绍如何使用 Grafana Tempo 收集和可视化 Gen AI 示例项目的遥测数据。

## 架构概述

```
Gen AI 应用 → OpenTelemetry HTTP Exporter → Tempo → Grafana
```

## 项目结构

```
examples/tempo/
├── docker-compose.yml              # Docker Compose 配置
├── tempo-config.yml               # Tempo 配置文件
├── TEMPO_SETUP.md                 # 本设置指南
└── grafana/                       # Grafana 配置和仪表板
    ├── provisioning/
    │   ├── datasources/
    │   │   └── tempo-datasource.yml
    │   └── dashboards/
    │       ├── tempo-dashboard.yml
    │       └── traces-dashboard.json
    └── data/                      # Grafana 数据存储
```

## 快速开始

### 1. 启动监控栈

```bash
# 启动 Tempo 和 Grafana
docker-compose up -d

# 查看启动日志
docker-compose logs -f tempo
```

### 2. 验证服务状态

```bash
# 检查服务是否运行
curl http://localhost:3200/tempo/api/v1/health

# 检查 Grafana 登录
curl -I http://localhost:3000
```

## 配置详解

### Tempo 配置 (`tempo-config.yml`)

```yaml
auth_enabled: false

server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        http:
          endpoint: "0.0.0.0:4318"
        grpc:
          endpoint: "0.0.0.0:4317"

storage:
  trace:
    backend: local
    local:
      path: /tmp/tempo/traces

compactor:
  working_directory: /tmp/tempo/blocks
  block_retention: 720h

query:
  frontend_service:
    template: "http://localhost:3200"
```

### 关键配置说明

- **OTLP接收器**: 支持 HTTP 和 gRPC 协议
- **本地存储**: 将追踪数据存储在本地文件系统
- **数据保留**: 30 天的追踪数据保留
- **查询服务**: 提供追踪查询和探索功能

### Grafana 配置

#### 数据源配置

```yaml
apiVersion: 1
datasources:
  - name: Tempo
    type: tempo
    access: proxy
    orgId: 1
    url: http://tempo:3200
    jsonData:
      nodeGraph:
        enabled: true
      traceQuery:
        enabled: true
      maxOutcomes: 500
      serviceMap:
        enabled: true
```

#### 预配置仪表板

项目包含以下仪表板：

1. **Trace 处理速率**: 显示每秒处理的追踪数量
2. **Trace 探索器**: 详细的追踪查看和过滤
3. **服务映射**: 服务间依赖关系可视化

## 使用示例

### 基本追踪查看

1. 访问 http://localhost:3000 (admin/admin)
2. 选择 "Gen AI Traces Dashboard"
3. 查看 "Trace Explorer" 面板
4. 可以根据服务名、追踪名等过滤追踪数据

### 高级查询

```sql
# 查询特定服务的追踪
service_name="gen-ai-example"

# 查询特定时间段的错误追踪
http.response.status_code >= 400

# 查询执行时间超过1秒的追踪
duration > 1s
```

## 生产环境部署

### 1. 使用分布式存储

```yaml
storage:
  trace:
    backend: s3
    s3:
      bucket: my-traces
      endpoint: s3.amazonaws.com
      access_key: ${AWS_ACCESS_KEY}
      secret_key: ${AWS_SECRET_KEY}
```

### 2. 添加认证

```yaml
auth_enabled: true
server:
  http_listen_port: 3200

auth:
  type: oidc
  oidc:
    client_id: tempo
    client_secret: ${OIDC_CLIENT_SECRET}
    issuer_url: ${OIDC_ISSUER_URL}
```

### 3. 使用高可用配置

```yaml
compactor:
  compactor_ring:
    addr: "127.0.0.1:9094"
    kvstore:
      store: memberlist
      prefix: "compactor-ring"
```

## 性能调优

### Tempo 性能优化

```yaml
compactor:
  max_block_bytes: 1073741824  # 1GB
  max_block_duration: 1h
  max_compaction_parallelism: 4

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          max_recv_msg_size: 4194304  # 4MB
          max_send_msg_size: 4194304
```

### 数据保留策略

```yaml
compactor:
  working_directory: /tmp/tempo/blocks
  block_retention: 168h  # 7 天
  compacted_block_retention: 336h  # 14 天
```

## 故障排除

### 常见问题

1. **追踪数据不显示**
   - 检查 Tempo 是否正常启动
   - 验证 OTLP 连接
   - 查查 Tempo 日志

2. **Grafana 连接失败**
   - 检查 Tempo 数据源配置
   - 验证网络连接
   - 检查防火墙设置

3. **性能问题**
   - 检查资源配置
   - 调整块大小和保留策略
   - 监控内存使用情况

### 调试命令

```bash
# 查看 Tempo 健康状态
curl http://localhost:3200/tempo/api/v1/health

# 查看块信息
curl http://localhost:3200/tempo/api/v1/blocks

# 查看 Tail API
curl http://localhost:3200/tempo/api/v1/traces/search?limit=10

# 查看 Tempo 日志
docker-compose logs tempo --tail=50
```

## 扩展功能

### 添加告警

在 Grafana 中配置以下告警：

1. **追踪处理速率下降**
2. **错误率超过阈值**
3. **延迟异常增长**
4. **服务不可用**

### 集成其他服务

```yaml
# 集成 Prometheus
query:
  frontend_service:
    template: |
      http://prometheus:9090

# 集成 Loki
backend: trace
trace:
  storage:
    trace:
      backend: loki
      loki:
        http:
          url: "http://loki:3100"
```

这个设置指南提供了完整的 Tempo 部署和使用说明，帮助您建立强大的追踪监控基础设施。