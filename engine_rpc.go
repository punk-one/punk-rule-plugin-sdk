package sdk

import (
	"fmt"
	"net/rpc"
)

// EngineRPC 引擎端 RPC 接口（v1.2）
// 插件通过此接口调用引擎能力：发布事件、记录日志、指标、Ack
type EngineRPC interface {
	Emit(e Event) error
	EmitWithTargets(e Event, toNodeIDs []string) error
	Log(level LogLevel, msg string, fields map[string]interface{})
	LogBatch(level LogLevel, messages []string, fields []map[string]interface{})
	IncCounter(name string, labels map[string]string)
	Observe(name string, value float64, labels map[string]string)
	Ack(eventID string) error
	PublishAck(ack AckMessage) error
	EmitBatch(events []Event) error
	GetState(key string) ([]byte, error)
	SetState(key string, value []byte) error
	DeleteState(key string) error
}

// EngineRPCClient 引擎端 RPC 客户端包装器
type EngineRPCClient struct {
	client *rpc.Client
}

func NewEngineRPCClient(client *rpc.Client) EngineRPC {
	return &EngineRPCClient{client: client}
}

func (c *EngineRPCClient) Emit(e Event) error {
	return c.EmitWithTargets(e, nil)
}

func (c *EngineRPCClient) EmitWithTargets(e Event, toNodeIDs []string) error {
	var reply struct{}
	eventJSON, err := EncodeEvent(e)
	if err != nil {
		return err
	}
	return c.client.Call("Engine.EmitRPC", struct {
		EventJSON []byte
		ToNodeIDs []string
	}{EventJSON: eventJSON, ToNodeIDs: toNodeIDs}, &reply)
}

func (c *EngineRPCClient) Log(level LogLevel, msg string, fields map[string]interface{}) {
	var reply struct{}
	_ = c.client.Call("Engine.LogRPC", struct {
		Level  LogLevel
		Msg    string
		Fields map[string]interface{}
	}{Level: level, Msg: msg, Fields: fields}, &reply)
}

func (c *EngineRPCClient) LogBatch(level LogLevel, messages []string, fields []map[string]interface{}) {
	var reply struct{}
	_ = c.client.Call("Engine.LogBatchRPC", &LogBatchArgs{
		Level:    level,
		Messages: messages,
		Fields:   fields,
	}, &reply)
}

func (c *EngineRPCClient) IncCounter(name string, labels map[string]string) {
	var reply struct{}
	_ = c.client.Call("Engine.IncCounterRPC", &MetricArgs{
		Name:   name,
		Labels: labels,
	}, &reply)
}

func (c *EngineRPCClient) Observe(name string, value float64, labels map[string]string) {
	var reply struct{}
	_ = c.client.Call("Engine.ObserveRPC", &MetricArgs{
		Name:   name,
		Value:  value,
		Labels: labels,
	}, &reply)
}

func (c *EngineRPCClient) Ack(eventID string) error {
	var reply struct{}
	return c.client.Call("Engine.AckRPC", &AckArgs{EventID: eventID}, &reply)
}

func (c *EngineRPCClient) PublishAck(ack AckMessage) error {
	var reply struct{}
	return c.client.Call("Engine.PublishAckRPC", &PublishAckArgs{Ack: ack}, &reply)
}

func (c *EngineRPCClient) EmitBatch(events []Event) error {
	var reply struct{}
	eventsJSON := make([][]byte, len(events))
	for i, e := range events {
		data, err := EncodeEvent(e)
		if err != nil {
			return fmt.Errorf("failed to marshal event %d: %w", i, err)
		}
		eventsJSON[i] = data
	}
	return c.client.Call("Engine.EmitBatchRPC", &EmitBatchArgs{
		EventsJSON: eventsJSON,
		ToNodeIDs:  nil,
	}, &reply)
}

func (c *EngineRPCClient) GetState(key string) ([]byte, error) {
	var reply StateGetReply
	if err := c.client.Call("Engine.GetStateRPC", &StateKeyArgs{Key: key}, &reply); err != nil {
		return nil, err
	}
	if !reply.Found {
		return nil, nil
	}
	return reply.Value, nil
}

func (c *EngineRPCClient) SetState(key string, value []byte) error {
	var reply struct{}
	return c.client.Call("Engine.SetStateRPC", &StateSetArgs{
		Key:   key,
		Value: value,
	}, &reply)
}

func (c *EngineRPCClient) DeleteState(key string) error {
	var reply struct{}
	return c.client.Call("Engine.DeleteStateRPC", &StateKeyArgs{Key: key}, &reply)
}

// RPC Args Structs

type LogBatchArgs struct {
	Level    LogLevel
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
	Ack AckMessage
}

type StateKeyArgs struct {
	Key string
}

type StateSetArgs struct {
	Key   string
	Value []byte
}

type StateGetReply struct {
	Value []byte
	Found bool
}

// Deprecated: EdgeSubject v1.3.1 已废弃
func EdgeSubject(ruleID, fromNodeID, toNodeID string) string {
	return fmt.Sprintf("rule.edge.%s.%s.%s", ruleID, fromNodeID, toNodeID)
}
