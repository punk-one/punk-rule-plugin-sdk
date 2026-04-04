# Punk Rule Engine Plugin SDK

Punk Rule Engine Plugin SDK 是一套高性能的 **"DAG 节点执行 SDK + Connector 资源插件 SDK"**，专为实时数据处理、边缘计算和云端规则编排而设计。

## 核心特性

- 🚀 **高性能 Batch RPC**: 支持批量事件处理和日志缓冲，大幅降低进程间通信开销。
- 🛠️ **统一接口**: Source、Processor、Sink 遵循同一套极简接口，易于开发和扩展。
- 📦 **进程隔离**: 基于 HashiCorp `go-plugin`，支持独立进程运行，确保系统稳定性。
- 📊 **多维观测**: 内置日志自聚合、自定义指标监控及全链路跟踪支持。
- 💾 **状态管理**: 专为有状态插件设计的 `StateManager`，支持状态持久化与迁移。
- ❤️ **健康上报**: `RuntimeContext.Health()` 成为一等公民，支持状态、恢复态和自动心跳。
- 🔌 **资源插件化**: 支持 `connector` 插件、共享连接、资源状态广播与 connector-aware Source/Sink。

## 安装

```bash
go get github.com/punk-one/punk-rule-plugin-sdk@latest
```

未发布 tag 时，`@latest` 会解析到默认分支最新提交对应的 pseudo-version。

发布正式版本后，推荐显式指定 tag：

```bash
go get github.com/punk-one/punk-rule-plugin-sdk@v1.7.1
```

## 目录结构

当前目录已按“稳定公共 API / 内部运行时 / 文档治理 / 示例测试”整理：

- `context.go`: `RuntimeContext`、`Logger`、`Emitter`、`Metrics`、`AckMessage` 及相关默认/桥接实现。
- `plugin.go`: `Plugin` 核心接口、`BasePlugin`、插件元信息类型与 go-plugin 兼容入口。
- `connector.go`: `ConnectorPlugin`、connector go-plugin 兼容入口。
- `serve.go`: `sdk.Serve(...)` 启动入口与 `ServeOptions`。
- `engine.go`: 引擎侧 RPC 接口与客户端桥接。
- `event.go`: 统一事件模型 `Event`、编码与 trace 辅助函数。
- `errors.go`: 结构化插件错误模型与错误分类辅助函数。
- `health.go`: `HealthReporter`、健康状态模型与上报配置。
- `state.go`: `StateManager`、状态存储抽象与默认实现。
- `legacy_compat.go`: 旧 engine / go-plugin / net/rpc DTO 的兼容导出层。
- `internal/healthruntime/`: 健康上报的异步、抑制与心跳运行时。
- `internal/runtime/`: 运行期上下文、logger、emitter、metrics、state 等内部实现。
- `internal/transport/`: `go-plugin`、RPC DTO、握手、跨进程调用等内部实现。
- `internal/core/`: 稳定模型、接口、事件、错误等核心定义。
- `schema/`: 规则流与消息契约相关 schema 类型。
- `examples/`: 可运行的 Source / Processor 示例。
- `tests/`: 面向公开 API 的黑盒测试与 API 基线测试。
- `docs/`: 接入文档，例如 `HealthReporter` 规范。
- `CHANGELOG.md`: 版本演进与重要变更记录。
- `public_api.json`: 机器可校验的公开 API 基线。
- `PUBLIC_API.md`: 从 API 基线生成的人类可读文档。

### 根目录放什么，`docs/` 放什么

- 根目录只保留三类内容：稳定公共 API 源码、仓库治理文件、发布入口文件。
- `README.md`、`CHANGELOG.md`、`PUBLIC_API.md`、`public_api.json` 保留在根目录，不放入 `docs/`：
  - `README.md` 是仓库首页入口；
  - `CHANGELOG.md` 是发布/版本记录；
  - `PUBLIC_API.md` 与 `public_api.json` 属于 API 治理基线，同时被测试与工具链使用。
- `docs/` 只放专题说明、接入文档、设计分析和模板文档。

## 快速开始

### 1. 实现一个转换插件 (Processor)

```go
package main

import (
    "github.com/punk-one/punk-rule-plugin-sdk"
)

type MyProcessor struct {
    sdk.BasePlugin
    ctx sdk.RuntimeContext
}

func (p *MyProcessor) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "my-processor",
        Name:    "My Processor",
        Version: "1.7.1",
        Type:    sdk.PluginTypeProcessor,
        Capabilities: sdk.PluginCapabilities{
            SupportBatch: true, // 声明支持批量处理
        },
    }
}

func (p *MyProcessor) Start(ctx sdk.RuntimeContext) error {
    p.ctx = ctx
    _ = ctx.Health().Healthy("processor started", map[string]string{
        "component": "my-processor",
    })
    return nil
}

// 单个事件处理
func (p *MyProcessor) OnEvent(e sdk.Event) error {
    e.Payload["processed"] = true
    return p.ctx.Emitter().Publish(e)
}

// 批量事件处理
func (p *MyProcessor) OnEvents(events []sdk.Event) error {
    for i := range events {
        events[i].Payload["processed"] = true
    }
    return p.ctx.Emitter().EmitBatch(events)
}

func main() {
    sdk.Serve(&MyProcessor{}, sdk.ServeOptions{})
}
```

更多可运行示例见：

- `examples/processor_basic/main.go`
- `examples/source_basic/main.go`
- `examples/health_basic/main.go`

其中 `examples/health_basic/main.go` 展示了如何在连接失败、运行中降级和周期性活跃场景下分别上报 `Unhealthy`、`Degraded` 和 `Heartbeat`。

### 3. 实现一个资源插件 (Connector)

```go
package main

import sdk "github.com/punk-one/punk-rule-plugin-sdk"

type MyConnector struct{}

func (c *MyConnector) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "connect-demo",
        Name:    "Demo Connector",
        Version: "1.7.1",
        Type:    sdk.PluginTypeConnector,
        Capabilities: sdk.PluginCapabilities{
            ConnectorDescriptor: &sdk.ConnectorDescriptor{
                Family:       "demo",
                Label:        "Demo Connector",
                Capabilities: []string{"read", "write", "shared_duplex"},
            },
        },
    }
}

func (c *MyConnector) CreateResource(resource sdk.ConnectorResource) (string, error) {
    return "demo-resource", nil
}

func (c *MyConnector) DestroyResource(providerHandle string) error { return nil }

func (c *MyConnector) Execute(providerHandle string, req sdk.ConnectorRequest) (sdk.ConnectorResponse, error) {
    return sdk.ConnectorResponse{}, nil
}

func (c *MyConnector) Probe(providerHandle string, req sdk.ConnectorRequest) (sdk.ResourceStatusEvent, error) {
    return sdk.ResourceStatusEvent{
        Status: sdk.ResourceStatusHealthy,
    }, nil
}

func (c *MyConnector) Stop() error { return nil }

func main() {
    sdk.ServeConnector(&MyConnector{})
}
```

### 2. 推荐插件模板

推荐从 `docs/plugin_template.md` 开始复制骨架，再按需裁剪为 Source / Processor / Sink。

## 核心 API

### Plugin 接口

```go
type Plugin interface {
    Info() PluginInfo
    Init(cfg PluginConfig) error
    Start(ctx RuntimeContext) error
    OnEvent(e Event) error         // 单点触发
    OnEvents(events []Event) error  // 批处理触发
    Stop() error
}
```

推荐插件直接嵌入：

```go
type MyPlugin struct {
    sdk.BasePlugin
}
```

这样只需要覆写实际会使用的方法；未实现的方法会返回结构化错误，便于引擎判定是否需要回退或中止。

如果插件希望覆盖默认健康心跳参数，也可以在启动入口传入：

```go
sdk.Serve(&MyPlugin{}, sdk.ServeOptions{
    Health: sdk.HealthOptions{
        HeartbeatInterval: 10 * time.Second,
        MaxSilencePeriod:  1 * time.Minute,
        QueueCapacity:     64,
    },
})
```

这里的 `ServeOptions` 是**插件进程启动期配置**，当前只暴露 `Health` 一个字段。

- `ServeOptions.Health` 表示插件进程级默认健康配置；
- 引擎在真正启动某个规则节点时，还可以通过 `StartArgs.Health` 传入本次启动的健康配置；
- SDK 会先读取启动参数，再回退到 `ServeOptions.Health` 作为默认值。

之所以不把这类配置直接塞进插件实现 `impl`，是因为它属于 **runtime/host 启动参数**，不是插件业务对象本身的属性。

### RuntimeContext 接口

```go
type RuntimeContext interface {
    RuleID() string
    NodeID() string
    Logger() Logger   // 缓冲式日志系统
    Emitter() Emitter // 高性能事件发射器
    Metrics() Metrics // 监控指标
    Health() HealthReporter // 节点健康与心跳
    Connector() ConnectorClient // 资源读写桥接
    ResourceEvents() <-chan ResourceStatusEvent // 资源状态事件
}
```

### ConnectorPlugin 接口

```go
type ConnectorPlugin interface {
    Info() PluginInfo
    CreateResource(resource ConnectorResource) (string, error)
    DestroyResource(providerHandle string) error
    Execute(providerHandle string, req ConnectorRequest) (ConnectorResponse, error)
    Probe(providerHandle string, req ConnectorRequest) (ResourceStatusEvent, error)
    Stop() error
}
```

### ConnectorResource 策略 helper

```go
resource := sdk.ConnectorResource{
    ID:       "plc-01",
    PluginID: "connect-s7",
}

_ = resource.SetHealthPolicy(sdk.HealthPolicy{
    Enabled:    true,
    IntervalMS: 1000,
    TimeoutMS:  500,
    Target:     "health",
    Priority:   "high",
})

### Connector 契约声明

对于 connector-aware 2.x 架构，建议插件在 `PluginCapabilities` 中显式声明 connector family 和绑定要求：

```go
func (s *S7Source) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "source-s7",
        Name:    "S7 Source",
        Version: "1.1.0",
        Type:    sdk.PluginTypeSource,
        Capabilities: sdk.PluginCapabilities{
            ConfigSchema: sourceSchema,
            ConnectorBinding: &sdk.ConnectorBindingSpec{
                Required:             true,
                AcceptedFamilies:     []string{"s7"},
                RequiredCapabilities: []string{"read"},
            },
        },
    }
}

func (c *S7Connector) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:      "connect-s7",
        Name:    "Connector Siemens S7",
        Version: "1.1.0",
        Type:    sdk.PluginTypeConnector,
        Capabilities: sdk.PluginCapabilities{
            ConfigSchema: connectorSchema,
            ConnectorDescriptor: &sdk.ConnectorDescriptor{
                Family:        "s7",
                Label:         "Siemens S7",
                Capabilities:  []string{"read", "write", "shared_duplex", "split_rw"},
                DisplayFields: []string{"host", "port", "rack", "slot"},
            },
        },
    }
}
```

推荐约定：

- `ConnectorDescriptor` 由 `connect-*` 插件声明自身的 family、可用能力和连接摘要字段。
- `ConnectorBindingSpec` 由 `source-*` / `sink-*` 插件声明其需要的 connector family 和 capability。
- 引擎与 UI 应优先依据这两个结构做兼容性过滤和强校验，而不是依赖字符串约定。

_ = resource.SetQuotaPolicy(sdk.QuotaPolicy{
    MaxInflightRequests: 16,
    MaxInflightWeight:   32,
    DefaultUsageWeight:  1,
})
```

### HealthReporter 接口

```go
type HealthReporter interface {
    Report(report HealthReport) error
    Healthy(msg string, details map[string]string) error
    Degraded(code, msg string, details map[string]string) error
    Recovering(code, msg string, details map[string]string) error
    Unhealthy(code, msg string, details map[string]string) error
    Heartbeat(msg string, details map[string]string) error
}
```

更多设计与接入约定见 `docs/health_reporter.md:1`。

### 结构化错误

```go
err := sdk.NewRetryableError("temporary PLC timeout")
if sdk.IsRetryable(err) {
    // engine can retry
}
```

SDK 会在 RPC 层保留错误语义，避免 `retryable`、`fatal`、`batch_unsupported` 等信息在跨进程调用时丢失。

### Emitter 接口

```go
type Emitter interface {
    Publish(event Event) error           // 发布单个事件
    EmitBatch(events []Event) error      // 批量发布
    EmitTo(label string, event Event) error // 定向路由发布
    Ack(eventID string) error            // 消息确认
}
```

### legacy_compat.go 的定位

`legacy_compat.go` 不是业务插件的推荐入口，而是根包保留的一层**兼容导出层**，主要用于：

- 旧引擎侧加载链路；
- `go-plugin` / `net/rpc` 相关 DTO 的历史公开符号；
- 官方插件或测试中仍直接引用的兼容类型。

例如 `StartArgs`、`ReceiveEventsArgs`、`LogBatchArgs` 这类类型，都属于跨进程协议兼容层，而不是业务插件开发首选 API。

如果未来决定：

- 同步升级引擎；
- 不再对外暴露这批底层 RPC DTO；
- 接受一次公开 API breaking change；

则可以考虑缩减甚至删除 `legacy_compat.go`。在那之前，建议继续保留它作为稳定 API 与兼容 API 的隔离层。

## 配置 Schema (ConfigSchema)

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
        Version: "1.7.0",
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
- **默认实现**: 推荐嵌入 `BasePlugin`，避免为不需要的方法编写样板代码。
- **批量回退**: 当插件未声明 `SupportBatch`，或 `OnEvents` 返回 `ErrBatchNotSupported` 时，SDK 会自动回退为逐条调用 `OnEvent`。
- **配置校验**: 强烈建议在 `PluginInfo` 中定义 `ConfigSchema`，引擎会在部署时据此进行严格的 JSON Schema 校验。
- **Schema 规范**: 只使用标准 JSON Schema 字段，不要使用 `ui:*` 扩展字段。
- **健康码规范**: 推荐使用 `SYS_*`、`NET_*`、`BIZ_*` 三类前缀，便于引擎聚合告警与统计。
- **API 治理**: 变更公开符号后，执行 `go generate ./...` 刷新 `PUBLIC_API.md`，并同步更新 `public_api.json` 基线。
- **模板起点**: 新插件建议从 `docs/plugin_template.md` 或 `examples/` 目录中的最小示例开始。
- **项目约定**: 根目录只保留稳定公共 API；示例放 `examples/`，测试放 `tests/`，补充文档放 `docs/`。

## 许可

Punk Rule Engine Plugin SDK 采用 Apache 2.0 许可证。
