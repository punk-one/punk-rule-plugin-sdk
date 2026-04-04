# Changelog

本文档记录 `punk-rule-plugin-sdk` 自 `v1.5.x` 以来的重要公共 API 与工程结构演进。

## v1.7.1 - 2026-04-04

### Added

- 新增 connector 契约公开类型：
  - `ConnectorDescriptor`
  - `ConnectorBindingSpec`

### Changed

- `PluginCapabilities` 新增两个可选字段：
  - `ConnectorDescriptor *ConnectorDescriptor`
  - `ConnectorBinding *ConnectorBindingSpec`
- README 已补充 connector family / capability / binding 的正式声明示例。
- `PUBLIC_API` 基线已扩展到 connector compatibility contract。

### Validation

- `go generate ./...`
- `go test ./...`

## v1.7.0 - 2026-04-04

### Added

- 新增 connector 插件稳定入口：
  - `PluginTypeConnector`
  - `ConnectorPlugin`
  - `ConnectorClient`
  - `ServeConnector(...)`
- 新增资源运行时能力：
  - `RuntimeContext.Connector()`
  - `RuntimeContext.ResourceEvents()`
- 新增 connector 公开模型：
  - `ConnectorResource`
  - `ConnectorRequest`
  - `ConnectorResponse`
  - `ResourceStatusEvent`
  - `ProviderPolicy`
  - `WarmupPolicy`
- 新增正式策略类型：
  - `HealthPolicy`
  - `QuotaPolicy`
- 新增 `ConnectorResource` 策略 helper：
  - `SetHealthPolicy(...)`
  - `GetHealthPolicy()`
  - `SetQuotaPolicy(...)`
  - `GetQuotaPolicy()`
- 新增 connector 公开 API 黑盒测试：`tests/connector_resource_test.go`

### Changed

- `PUBLIC_API` 基线已扩展到 connector-aware 2.x 运行时。
- `PluginRuntimeContext` 现在会为外部插件桥接 connector 调用与资源状态事件流。
- README、公开 API 基线与测试已对齐 connector 插件开发模式。

### Validation

- `go test ./tests ./internal/runtime` 已通过。

## v1.6.0 - 2026-03-27

### Added

- 新增 `BasePlugin`，为插件作者提供无副作用生命周期默认实现。
- 新增结构化错误模型：`PluginError`、`ErrorKind`、`ErrNotImplemented`、`ErrBatchNotSupported`。
- 新增事件工具：`Event.Clone()`、`GetTraceID(...)`、`SetTraceID(...)`。
- 新增 API 治理文件：`public_api.json`、`PUBLIC_API.md`、`internal/tools/public_api_export.go`。
- 新增健康接入文档：`docs/health_reporter.md`。
- 新增可运行示例：`examples/health_basic/main.go`。
- 新增维护入口：`doc.go`，支持 `go generate ./...` 刷新公开 API 文档。

### Changed

- `PluginRPC` 现在支持结构化错误透传，不再丢失重试/致命/批量不支持语义。
- 批量分发链路支持自动回退：未声明 `SupportBatch` 或运行时返回 `ErrBatchNotSupported` 时，自动退回逐条 `OnEvent(...)`。
- `ServeOptions.Health` 现在会参与插件启动时的默认健康配置合并。
- `Stop()` 失败时的健康上报消息改为原始错误文本，不再暴露内部 RPC 编码串。
- README、示例、测试和公开 API 基线已对齐。

### Refactored

- 根目录职责收敛为“稳定公共 API”；示例统一放入 `examples/`，黑盒测试统一放入 `tests/`。
- `runtime_context`、`emitter`、`metrics`、`engine_rpc` 已按“接口 / 实现 / RPC 适配”拆分为更清晰的文件。

### Removed

- 从根包移除实现细节导出：`PluginRuntimeContext`、`NewPluginRuntimeContext`、`EngineRPCEmitter`、`EngineRPCLogger`、`EngineRPCMetrics`、`NoOpEmitter`、`StateManagerImpl`。
- SDK 自身测试改为优先依赖稳定接口，运行时实现测试下沉到 `internal/runtime`。

### Validation

- `go test ./...` 已通过。
