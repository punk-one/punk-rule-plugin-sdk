# Public API Surface

Generated from `public_api.json` via `go run ./internal/tools/public_api_export.go`.

## Machine-readable Source

```json
{
  "symbols": [
    "AckArgs",
    "AckMessage",
    "AckStatusFailed",
    "AckStatusSuccess",
    "BasePlugin",
    "BuiltinRuntimeContext",
    "DecodeEvent",
    "DefaultHealthOptions",
    "DefaultLogger",
    "DefaultMetrics",
    "EdgeSubject",
    "EmitBatchArgs",
    "Emitter",
    "EmitterAdapter",
    "EncodeEvent",
    "EngineBridgeRuntimeContext",
    "EngineRPC",
    "EngineRPCClient",
    "ErrBatchNotSupported",
    "ErrNotImplemented",
    "ErrorKind",
    "ErrorKindBatchUnsupported",
    "ErrorKindFatal",
    "ErrorKindNotImplemented",
    "ErrorKindRetryable",
    "ErrorKindSkippable",
    "Event",
    "EventMetadataTraceIDKey",
    "FatalError",
    "Field",
    "GetTraceID",
    "HandshakeConfig",
    "HealthDegraded",
    "HealthHealthy",
    "HealthHeartbeat",
    "HealthInitializing",
    "HealthKind",
    "HealthKindHeartbeat",
    "HealthKindState",
    "HealthOptions",
    "HealthRecovering",
    "HealthReport",
    "HealthReporter",
    "HealthStatus",
    "HealthStopped",
    "HealthUnhealthy",
    "HealthUnknown",
    "InfoReply",
    "InitArgs",
    "InitReply",
    "IsBatchNotSupported",
    "IsFatal",
    "IsRetryable",
    "IsSkippable",
    "LogBatchArgs",
    "LogLevel",
    "LogLevelDebug",
    "LogLevelError",
    "LogLevelInfo",
    "LogLevelWarn",
    "Logger",
    "MetricArgs",
    "Metrics",
    "NewEngineRPCClient",
    "NewEvent",
    "NewFatalError",
    "NewRPCEmitterFromEngine",
    "NewRPCLoggerFromEngine",
    "NewRPCMetricsFromEngine",
    "NewRetryableError",
    "NewSkippableError",
    "NewStateManager",
    "Plugin",
    "PluginCapabilities",
    "PluginConfig",
    "PluginError",
    "PluginInfo",
    "PluginRPC",
    "PluginRPCClient",
    "PluginRPCServer",
    "PluginType",
    "PluginTypeProcessor",
    "PluginTypeSink",
    "PluginTypeSource",
    "PluginTypeUtility",
    "PublishAckArgs",
    "ReceiveEventArgs",
    "ReceiveEventReply",
    "ReceiveEventsArgs",
    "ReceiveEventsReply",
    "ReportHealthArgs",
    "RetryableError",
    "RuntimeContext",
    "RuntimeContextImpl",
    "Serve",
    "ServeOptions",
    "SetTraceID",
    "SkippableError",
    "StartArgs",
    "StartReply",
    "StateGetReply",
    "StateKeyArgs",
    "StateManager",
    "StateSetArgs",
    "StateSetWithTTLArgs",
    "StateStore",
    "StatefulRuntimeContext",
    "StopReply"
  ]
}
```

## Symbols

- `AckArgs`
- `AckMessage`
- `AckStatusFailed`
- `AckStatusSuccess`
- `BasePlugin`
- `BuiltinRuntimeContext`
- `DecodeEvent`
- `DefaultHealthOptions`
- `DefaultLogger`
- `DefaultMetrics`
- `EdgeSubject`
- `EmitBatchArgs`
- `Emitter`
- `EmitterAdapter`
- `EncodeEvent`
- `EngineBridgeRuntimeContext`
- `EngineRPC`
- `EngineRPCClient`
- `ErrBatchNotSupported`
- `ErrNotImplemented`
- `ErrorKind`
- `ErrorKindBatchUnsupported`
- `ErrorKindFatal`
- `ErrorKindNotImplemented`
- `ErrorKindRetryable`
- `ErrorKindSkippable`
- `Event`
- `EventMetadataTraceIDKey`
- `FatalError`
- `Field`
- `GetTraceID`
- `HandshakeConfig`
- `HealthDegraded`
- `HealthHealthy`
- `HealthHeartbeat`
- `HealthInitializing`
- `HealthKind`
- `HealthKindHeartbeat`
- `HealthKindState`
- `HealthOptions`
- `HealthRecovering`
- `HealthReport`
- `HealthReporter`
- `HealthStatus`
- `HealthStopped`
- `HealthUnhealthy`
- `HealthUnknown`
- `InfoReply`
- `InitArgs`
- `InitReply`
- `IsBatchNotSupported`
- `IsFatal`
- `IsRetryable`
- `IsSkippable`
- `LogBatchArgs`
- `LogLevel`
- `LogLevelDebug`
- `LogLevelError`
- `LogLevelInfo`
- `LogLevelWarn`
- `Logger`
- `MetricArgs`
- `Metrics`
- `NewEngineRPCClient`
- `NewEvent`
- `NewFatalError`
- `NewRPCEmitterFromEngine`
- `NewRPCLoggerFromEngine`
- `NewRPCMetricsFromEngine`
- `NewRetryableError`
- `NewSkippableError`
- `NewStateManager`
- `Plugin`
- `PluginCapabilities`
- `PluginConfig`
- `PluginError`
- `PluginInfo`
- `PluginRPC`
- `PluginRPCClient`
- `PluginRPCServer`
- `PluginType`
- `PluginTypeProcessor`
- `PluginTypeSink`
- `PluginTypeSource`
- `PluginTypeUtility`
- `PublishAckArgs`
- `ReceiveEventArgs`
- `ReceiveEventReply`
- `ReceiveEventsArgs`
- `ReceiveEventsReply`
- `ReportHealthArgs`
- `RetryableError`
- `RuntimeContext`
- `RuntimeContextImpl`
- `Serve`
- `ServeOptions`
- `SetTraceID`
- `SkippableError`
- `StartArgs`
- `StartReply`
- `StateGetReply`
- `StateKeyArgs`
- `StateManager`
- `StateSetArgs`
- `StateSetWithTTLArgs`
- `StateStore`
- `StatefulRuntimeContext`
- `StopReply`
