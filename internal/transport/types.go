package transport

import "github.com/punk-one/punk-rule-plugin-sdk/internal/core"

type LogBatchArgs struct {
	Level    core.LogLevel
	Messages []string
	Fields   []map[string]interface{}
}

type EmitBatchArgs struct {
	EventsJSON [][]byte
	ToNodeIDs  []string
}

type MetricArgs struct {
	Name   string
	Value  float64
	Labels map[string]string
}

type AckArgs struct {
	EventID string
}

type PublishAckArgs struct {
	Ack core.AckMessage
}

type StateKeyArgs struct {
	Key string
}

type StateSetArgs struct {
	Key   string
	Value []byte
}

type StateSetWithTTLArgs struct {
	Key      string
	Value    []byte
	TTLNanos int64
}

type StateGetReply struct {
	Value []byte
	Found bool
}

type InfoReply struct{ Info core.PluginInfo }
type InitArgs struct{ Config core.PluginConfig }
type InitReply struct{ Error string }

type StartArgs struct {
	RuleID string
	NodeID string
	Health core.HealthOptions
}

type StartReply struct{ Error string }
type ReceiveEventArgs struct{ EventJSON []byte }
type ReceiveEventReply struct{ Error string }
type ReceiveEventsArgs struct{ EventsJSON []byte }
type ReceiveEventsReply struct{ Error string }
type StopReply struct{ Error string }
