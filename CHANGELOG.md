# Changelog

本文档记录 `punk-rule-plugin-sdk` 自 `v1.5.x` 以来的重要公共 API 与工程结构演进。

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
