package core

import "encoding/json"

type PluginType string

const (
	PluginTypeSource    PluginType = "source"
	PluginTypeProcessor PluginType = "processor"
	PluginTypeSink      PluginType = "sink"
	PluginTypeUtility   PluginType = "utility"
)

type PluginInfo struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Type         PluginType         `json:"type"`
	Description  string             `json:"description"`
	Author       string             `json:"author"`
	Category     string             `json:"category"`
	Capabilities PluginCapabilities `json:"capabilities"`
}

type PluginCapabilities struct {
	InputPorts   int    `json:"input_ports"`
	OutputPorts  int    `json:"output_ports"`
	SupportBatch bool   `json:"support_batch"`
	SupportAck   bool   `json:"support_ack"`
	ConfigSchema string `json:"config_schema"`
	Stateful     bool   `json:"stateful"`
}

type PluginConfig struct {
	Raw    json.RawMessage `json:"raw"`
	RuleID string          `json:"rule_id,omitempty"`
	NodeID string          `json:"node_id,omitempty"`
	Health HealthOptions   `json:"health,omitempty"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type Field struct {
	Key   string
	Value interface{}
}
