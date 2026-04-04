package sdk

import (
	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
	transport "github.com/punk-one/punk-rule-plugin-sdk/internal/transport"
)

type Plugin = core.Plugin
type BasePlugin = core.BasePlugin

type PluginType = core.PluginType

const (
	PluginTypeSource    = core.PluginTypeSource
	PluginTypeProcessor = core.PluginTypeProcessor
	PluginTypeSink      = core.PluginTypeSink
	PluginTypeUtility   = core.PluginTypeUtility
	PluginTypeConnector = core.PluginTypeConnector
)

type PluginInfo = core.PluginInfo
type PluginCapabilities = core.PluginCapabilities
type PluginConfig = core.PluginConfig
type ConnectorClient = core.ConnectorClient
type ConnectorResource = core.ConnectorResource
type ConnectorRequest = core.ConnectorRequest
type ConnectorRequestKind = core.ConnectorRequestKind
type ConnectorResponse = core.ConnectorResponse
type HealthPolicy = core.HealthPolicy
type ProviderPolicy = core.ProviderPolicy
type QuotaPolicy = core.QuotaPolicy
type ResourceStatus = core.ResourceStatus
type ResourceStatusEvent = core.ResourceStatusEvent
type WarmupPolicy = core.WarmupPolicy

const (
	WarmupPolicyLazy     = core.WarmupPolicyLazy
	WarmupPolicyPreload  = core.WarmupPolicyPreload
	WarmupPolicyAlwaysOn = core.WarmupPolicyAlwaysOn

	ResourceStatusUnknown      = core.ResourceStatusUnknown
	ResourceStatusHealthy      = core.ResourceStatusHealthy
	ResourceStatusDegraded     = core.ResourceStatusDegraded
	ResourceStatusCongested    = core.ResourceStatusCongested
	ResourceStatusDisconnected = core.ResourceStatusDisconnected
	ResourceStatusReadonly     = core.ResourceStatusReadonly
	ResourceStatusBlocked      = core.ResourceStatusBlocked
	ResourceStatusRecovering   = core.ResourceStatusRecovering

	ConnectorRequestRead       = core.ConnectorRequestRead
	ConnectorRequestWrite      = core.ConnectorRequestWrite
	ConnectorRequestBatchRead  = core.ConnectorRequestBatchRead
	ConnectorRequestBatchWrite = core.ConnectorRequestBatchWrite
	ConnectorRequestMetadata   = core.ConnectorRequestMetadata
)

type LogLevel = core.LogLevel

const (
	LogLevelDebug = core.LogLevelDebug
	LogLevelInfo  = core.LogLevelInfo
	LogLevelWarn  = core.LogLevelWarn
	LogLevelError = core.LogLevelError
)

type Field = core.Field

// PluginRPC 实现了 go-plugin 的 Plugin 接口
type PluginRPC = transport.PluginRPC

// PluginRPCServer RPC 服务端
type PluginRPCServer = transport.PluginRPCServer

// PluginRPCClient RPC 客户端
type PluginRPCClient = transport.PluginRPCClient

// HandshakeConfig 是 go-plugin 的握手配置
// 确保插件和引擎使用的协议版本一致
var HandshakeConfig = transport.HandshakeConfig
