# Punk Rule Engine Plugin SDK v1.5.0

Punk Rule Engine Plugin SDK 是一套高性能的 **"DAG 节点执行 SDK"**，专为实时数据处理、边缘计算和云端规则编排而设计。

## 核心特性

- 🚀 **高性能 Batch RPC (v1.5.0)**: 支持批量事件处理和日志缓冲，大幅降低进程间通信开销。
- 🛠️ **统一接口**: Source、Processor、Sink 遵循同一套极简接口，易于开发和扩展。
- 📦 **进程隔离**: 基于 HashiCorp `go-plugin`，支持独立进程运行，确保系统稳定性。
- 📊 **多维观测**: 内置日志自聚合、自定义指标监控及全链路跟踪支持。
- 💾 **状态管理**: 专为有状态插件设计的 `StateManager`，支持状态持久化与迁移。

## 安装

```bash
go get github.com/punk-one/punk-rule-plugin-sdk@latest
```

未发布 tag 时，`@latest` 会解析到默认分支最新提交对应的 pseudo-version。

发布正式版本后，推荐显式指定 tag：

```bash
go get github.com/punk-one/punk-rule-plugin-sdk@v1.5.0
```

## 目录结构 (v1.5.0 重构版)

为了提高代码的可维护性和可读性，SDK 进行了结构化重构：

- `interfaces.go`: 定义所有核心接口 (`Plugin`, `RuntimeContext`, `Emitter`, `Logger`, etc.)。
- `types.go`: 定义基础数据类型和枚举 (`PluginInfo`, `LogLevel`, `Field`, etc.)。
- `event.go`: 定义统一事件模型 `Event` 及其编码辅助函数。
- `defaults.go`: 提供默认的 Logger、Metrics 和内建 Runtime 适配。
- `bridge.go`: RPC 桥接实现，包含高性能日志记录器和事件发射器。
- `plugin_rpc.go`: `go-plugin` 的 RPC 协议实现。
- `engine_rpc.go`: 定义引擎端提供的 RPC 服务接口。
- `handshake.go`: `go-plugin` 握手配置。

## 快速开始

### 1. 实现一个转换插件 (Processor)

```go
package main

import (
    "github.com/hashicorp/go-plugin"
    "github.com/punk-one/punk-rule-plugin-sdk"
)

type MyProcessor struct {
    ctx sdk.RuntimeContext
}

func (p *MyProcessor) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "my-processor",
        Name:    "My Processor",
        Version: "1.5.0",
        Type:    sdk.PluginTypeProcessor,
        Capabilities: sdk.PluginCapabilities{
            SupportBatch: true, // 声明支持批量处理
        },
    }
}

func (p *MyProcessor) Init(cfg sdk.PluginConfig) error {
    return nil // 解析配置
}

func (p *MyProcessor) Start(ctx sdk.RuntimeContext) error {
    p.ctx = ctx
    return nil
}

// 单个事件处理
func (p *MyProcessor) OnEvent(e sdk.Event) error {
    e.Payload["processed"] = true
    return p.ctx.Emitter().Publish(e)
}

// 批量事件处理 (v1.5.0 性能优化点)
func (p *MyProcessor) OnEvents(events []sdk.Event) error {
    for i := range events {
        events[i].Payload["processed"] = true
    }
    return p.ctx.Emitter().EmitBatch(events)
}

func (p *MyProcessor) Stop() error {
    return nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: sdk.HandshakeConfig,
        Plugins: map[string]plugin.Plugin{
            "plugin": &sdk.PluginRPC{Impl: &MyProcessor{}},
        },
    })
}
```

## 核心 API

### Plugin 接口 (v1.5.0)

```go
type Plugin interface {
    Info() PluginInfo
    Init(cfg PluginConfig) error
    Start(ctx RuntimeContext) error
    OnEvent(e Event) error         // 单点触发
    OnEvents(events []Event) error  // 批处理触发 (v1.5.0)
    Stop() error
}
```

### RuntimeContext 接口

```go
type RuntimeContext interface {
    RuleID() string
    NodeID() string
    Logger() Logger   // 缓冲式日志系统
    Emitter() Emitter // 高性能事件发射器
    Metrics() Metrics // 监控指标
}
```

### Emitter 接口

```go
type Emitter interface {
    Publish(event Event) error           // 发布单个事件
    EmitBatch(events []Event) error      // 批量发布 (v1.5.0 高性能)
    EmitTo(label string, event Event) error // 定向路由发布 (v1.3.1)
    Ack(eventID string) error            // 消息确认
}
```

## 配置 Schema (ConfigSchema) v1.5.0

### 规范要求

插件必须通过 `PluginInfo.Capabilities.ConfigSchema` 提供 **纯 JSON Schema**，用于定义插件配置的结构和约束。

**重要约束**：
- ✅ **只使用标准 JSON Schema 字段**：`type`、`properties`、`required`、`enum`、`default`、`title`、`description`、`format` 等
- ❌ **禁止使用 UI 扩展字段**：不得在 Schema 中使用 `ui:*` 字段（如 `ui:widget`、`ui:options`、`ui:collapsible` 等）
- ✅ **前端自动注入**：前端会根据 Schema 自动为 `input`、`process`、`state`、`output` 字段注入 `ui:options`，实现折叠面板等功能

### 示例

```go
func (p *MyProcessor) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "my-processor",
        Name:    "My Processor",
        Version: "1.5.0",
        Type:    sdk.PluginTypeProcessor,
        Capabilities: sdk.PluginCapabilities{
            ConfigSchema: `{
                "$schema": "https://json-schema.org/draft/2020-12/schema",
                "type": "object",
                "properties": {
                    "process": {
                        "type": "object",
                        "properties": {
                            "threshold": {
                                "type": "number",
                                "title": "Threshold",
                                "description": "Processing threshold value",
                                "default": 0.5
                            },
                            "filter_expr": {
                                "type": "string",
                                "title": "Filter Expression",
                                "description": "CEL expression to filter events",
                                "format": "cel"
                            }
                        },
                        "required": ["threshold"]
                    }
                }
            }`,
        },
    }
}
```

### Widget 自动选择规则

前端会根据 Schema 的 `type` 和 `format` 自动选择合适的 Widget：

| Schema 特征 | Widget |
|-----------|--------|
| `string` + `format: textarea` | Textarea |
| `string` + `format: json/cel/code` | Monaco Editor |
| `string` + `enum` | Select |
| `string` (default) | Text Input |
| `boolean` | Checkbox |
| `number`/`integer` | Number Input |
| `object` | ObjectFieldTemplate (折叠面板) |
| `array` | ArrayFieldTemplate |

### 自动折叠面板

以下字段名会自动启用折叠面板（前端自动注入 `ui:options`）：
- `input`
- `process`
- `state`
- `output`

插件作者无需在 Schema 中声明，前端会自动处理。

## 开发者提示

- **日志性能**: `ctx.Logger()` 已内置高性能异步缓冲，默认每 100 条或 2 秒自动刷新，无需担心高频日志阻塞 RPC。
- **并发安全**: `OnEvent` 可能并发调用（取决于引擎调度），请确保插件内部状态处理使用了互斥锁。
- **配置校验**: 强烈建议在 `PluginInfo` 中定义 `ConfigSchema`，引擎会在部署时据此进行严格的 JSON Schema 校验。
- **Schema 规范**: 遵循 v1.5.0 规范，只使用标准 JSON Schema 字段，不要使用 `ui:*` 扩展字段。

## 许可

Punk Rule Engine Plugin SDK 采用 Apache 2.0 许可证。
