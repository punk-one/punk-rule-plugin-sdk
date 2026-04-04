package core

import "encoding/json"

type PluginType string

const (
	PluginTypeSource    PluginType = "source"
	PluginTypeProcessor PluginType = "processor"
	PluginTypeSink      PluginType = "sink"
	PluginTypeUtility   PluginType = "utility"
	PluginTypeConnector PluginType = "connector"
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

type ConnectorDescriptor struct {
	Family        string   `json:"family"`
	Label         string   `json:"label,omitempty"`
	Capabilities  []string `json:"capabilities,omitempty"`
	DisplayFields []string `json:"display_fields,omitempty"`
}

type ConnectorBindingSpec struct {
	Required             bool     `json:"required,omitempty"`
	AcceptedFamilies     []string `json:"accepted_families,omitempty"`
	RequiredCapabilities []string `json:"required_capabilities,omitempty"`
}

type PluginCapabilities struct {
	InputPorts          int                   `json:"input_ports"`
	OutputPorts         int                   `json:"output_ports"`
	SupportBatch        bool                  `json:"support_batch"`
	SupportAck          bool                  `json:"support_ack"`
	ConfigSchema        string                `json:"config_schema"`
	Stateful            bool                  `json:"stateful"`
	ConnectorDescriptor *ConnectorDescriptor  `json:"connector_descriptor,omitempty"`
	ConnectorBinding    *ConnectorBindingSpec `json:"connector_binding,omitempty"`
}

type PluginConfig struct {
	Raw    json.RawMessage `json:"raw"`
	RuleID string          `json:"rule_id,omitempty"`
	NodeID string          `json:"node_id,omitempty"`
	Health HealthOptions   `json:"health,omitempty"`
}

type WarmupPolicy string

const (
	WarmupPolicyLazy     WarmupPolicy = "lazy"
	WarmupPolicyPreload  WarmupPolicy = "preload"
	WarmupPolicyAlwaysOn WarmupPolicy = "always_on"
)

type ResourceStatus string

const (
	ResourceStatusUnknown      ResourceStatus = "unknown"
	ResourceStatusHealthy      ResourceStatus = "healthy"
	ResourceStatusDegraded     ResourceStatus = "degraded"
	ResourceStatusCongested    ResourceStatus = "congested"
	ResourceStatusDisconnected ResourceStatus = "disconnected"
	ResourceStatusReadonly     ResourceStatus = "readonly"
	ResourceStatusBlocked      ResourceStatus = "blocked"
	ResourceStatusRecovering   ResourceStatus = "recovering"
)

type ProviderPolicy struct {
	MaxResourcesPerProvider int `json:"max_resources_per_provider,omitempty"`
	MaxInflightPerProvider  int `json:"max_inflight_per_provider,omitempty"`
	PreferredProviderCount  int `json:"preferred_provider_count,omitempty"`
}

type HealthPolicy struct {
	Enabled    bool   `json:"enabled,omitempty"`
	IntervalMS int    `json:"interval_ms,omitempty"`
	TimeoutMS  int    `json:"timeout_ms,omitempty"`
	Target     string `json:"target,omitempty"`
	Priority   string `json:"priority,omitempty"`
}

type QuotaPolicy struct {
	MaxInflightRequests int `json:"max_inflight_requests,omitempty"`
	MaxInflightWeight   int `json:"max_inflight_weight,omitempty"`
	DefaultUsageWeight  int `json:"default_usage_weight,omitempty"`
}

type ConnectorResource struct {
	ID             string            `json:"id"`
	Name           string            `json:"name,omitempty"`
	PluginID       string            `json:"plugin_id"`
	Connection     json.RawMessage   `json:"connection,omitempty"`
	SessionPolicy  string            `json:"session_policy,omitempty"`
	WarmupPolicy   WarmupPolicy      `json:"warmup_policy,omitempty"`
	HealthPolicy   json.RawMessage   `json:"health_policy,omitempty"`
	QuotaPolicy    json.RawMessage   `json:"quota_policy,omitempty"`
	ProviderPolicy ProviderPolicy    `json:"provider_policy,omitempty"`
	Enabled        bool              `json:"enabled"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

func (r *ConnectorResource) SetHealthPolicy(policy HealthPolicy) error {
	if r == nil {
		return nil
	}
	if isZeroHealthPolicy(policy) {
		r.HealthPolicy = nil
		return nil
	}
	raw, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	r.HealthPolicy = raw
	return nil
}

func (r ConnectorResource) GetHealthPolicy() (HealthPolicy, error) {
	if len(r.HealthPolicy) == 0 {
		return HealthPolicy{}, nil
	}
	var policy HealthPolicy
	if err := json.Unmarshal(r.HealthPolicy, &policy); err != nil {
		return HealthPolicy{}, err
	}
	return policy, nil
}

func (r *ConnectorResource) SetQuotaPolicy(policy QuotaPolicy) error {
	if r == nil {
		return nil
	}
	if isZeroQuotaPolicy(policy) {
		r.QuotaPolicy = nil
		return nil
	}
	raw, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	r.QuotaPolicy = raw
	return nil
}

func (r ConnectorResource) GetQuotaPolicy() (QuotaPolicy, error) {
	if len(r.QuotaPolicy) == 0 {
		return QuotaPolicy{}, nil
	}
	var policy QuotaPolicy
	if err := json.Unmarshal(r.QuotaPolicy, &policy); err != nil {
		return QuotaPolicy{}, err
	}
	return policy, nil
}

func isZeroHealthPolicy(policy HealthPolicy) bool {
	return !policy.Enabled &&
		policy.IntervalMS == 0 &&
		policy.TimeoutMS == 0 &&
		policy.Target == "" &&
		policy.Priority == ""
}

func isZeroQuotaPolicy(policy QuotaPolicy) bool {
	return policy.MaxInflightRequests == 0 &&
		policy.MaxInflightWeight == 0 &&
		policy.DefaultUsageWeight == 0
}

type ConnectorRequestKind string

const (
	ConnectorRequestRead       ConnectorRequestKind = "read"
	ConnectorRequestWrite      ConnectorRequestKind = "write"
	ConnectorRequestBatchRead  ConnectorRequestKind = "batch_read"
	ConnectorRequestBatchWrite ConnectorRequestKind = "batch_write"
	ConnectorRequestMetadata   ConnectorRequestKind = "metadata"
)

type ConnectorRequest struct {
	ResourceToken string               `json:"resource_token,omitempty"`
	Kind          ConnectorRequestKind `json:"kind"`
	Target        string               `json:"target,omitempty"`
	Payload       json.RawMessage      `json:"payload,omitempty"`
	Priority      string               `json:"priority,omitempty"`
	UsageWeight   int                  `json:"usage_weight,omitempty"`
	TimeoutMS     int                  `json:"timeout_ms,omitempty"`
}

type ConnectorResponse struct {
	Payload      json.RawMessage   `json:"payload,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Congested    bool              `json:"congested,omitempty"`
	RetryAfterMS int               `json:"retry_after_ms,omitempty"`
}

type ResourceStatusEvent struct {
	ResourceID         string            `json:"resource_id,omitempty"`
	ProviderPluginID   string            `json:"provider_plugin_id,omitempty"`
	ProviderInstanceID string            `json:"provider_instance_id,omitempty"`
	Status             ResourceStatus    `json:"status"`
	ReasonCode         string            `json:"reason_code,omitempty"`
	Message            string            `json:"message,omitempty"`
	Details            map[string]string `json:"details,omitempty"`
	ObservedAtUnixNano int64             `json:"observed_at_unix_nano,omitempty"`
	Sequence           uint64            `json:"sequence,omitempty"`
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
