package sdk

import (
	"encoding/json"
)

// PluginType 插件类型 (v1.1.0)
type PluginType string

const (
	PluginTypeSource    PluginType = "source"
	PluginTypeProcessor PluginType = "processor"
	PluginTypeSink      PluginType = "sink"
	PluginTypeUtility   PluginType = "utility"
)

// PluginInfo 插件元数据信息 (v1.1.0)
type PluginInfo struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Type         PluginType         `json:"type"`
	Description  string             `json:"description"`
	Author       string             `json:"author"`
	Category     string             `json:"category"` // v1.1.1: 插件分类
	Capabilities PluginCapabilities `json:"capabilities"`
}

// PluginCapabilities 插件能力声明
type PluginCapabilities struct {
	InputPorts   int    `json:"input_ports"`   // 输入端口数量
	OutputPorts  int    `json:"output_ports"`  // 输出端口数量
	SupportBatch bool   `json:"support_batch"` // 是否支持批量处理 (v1.5.0)
	SupportAck   bool   `json:"support_ack"`   // 是否支持确认机制
	ConfigSchema string `json:"config_schema"` // JSON Schema (强校验)
	Stateful     bool   `json:"stateful"`      // 是否是有状态插件
}

// PluginConfig 插件配置
type PluginConfig struct {
	Raw json.RawMessage `json:"raw"`
}

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}
