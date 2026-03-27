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
)

type PluginInfo = core.PluginInfo
type PluginCapabilities = core.PluginCapabilities
type PluginConfig = core.PluginConfig

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
