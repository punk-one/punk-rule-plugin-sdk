package transport

import (
	"errors"
	"fmt"
	"net/rpc"
	"time"

	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
	"github.com/punk-one/punk-rule-plugin-sdk/internal/runtime"
)

type EngineRPCClient struct {
	client *rpc.Client
}

func NewEngineRPCClient(client *rpc.Client) *EngineRPCClient {
	return &EngineRPCClient{client: client}
}

func (c *EngineRPCClient) Emit(e core.Event) error {
	return c.EmitWithTargets(e, nil)
}

func (c *EngineRPCClient) EmitWithTargets(e core.Event, toNodeIDs []string) error {
	var reply struct{}
	eventJSON, err := core.EncodeEvent(e)
	if err != nil {
		return err
	}
	return c.client.Call("Engine.EmitRPC", struct {
		EventJSON []byte
		ToNodeIDs []string
	}{EventJSON: eventJSON, ToNodeIDs: toNodeIDs}, &reply)
}

func (c *EngineRPCClient) Log(level core.LogLevel, msg string, fields map[string]interface{}) {
	var reply struct{}
	_ = c.client.Call("Engine.LogRPC", struct {
		Level  core.LogLevel
		Msg    string
		Fields map[string]interface{}
	}{Level: level, Msg: msg, Fields: fields}, &reply)
}

func (c *EngineRPCClient) LogBatch(level core.LogLevel, messages []string, fields []map[string]interface{}) {
	var reply struct{}
	_ = c.client.Call("Engine.LogBatchRPC", &LogBatchArgs{
		Level:    level,
		Messages: messages,
		Fields:   fields,
	}, &reply)
}

func (c *EngineRPCClient) IncCounter(name string, labels map[string]string) {
	var reply struct{}
	_ = c.client.Call("Engine.IncCounterRPC", &MetricArgs{Name: name, Labels: labels}, &reply)
}

func (c *EngineRPCClient) Observe(name string, value float64, labels map[string]string) {
	var reply struct{}
	_ = c.client.Call("Engine.ObserveRPC", &MetricArgs{Name: name, Value: value, Labels: labels}, &reply)
}

func (c *EngineRPCClient) ReportHealth(args core.ReportHealthArgs) error {
	var reply struct{}
	return c.client.Call("Engine.ReportHealthRPC", &args, &reply)
}

func (c *EngineRPCClient) Ack(eventID string) error {
	var reply struct{}
	return c.client.Call("Engine.AckRPC", &AckArgs{EventID: eventID}, &reply)
}

func (c *EngineRPCClient) PublishAck(ack core.AckMessage) error {
	var reply struct{}
	return c.client.Call("Engine.PublishAckRPC", &PublishAckArgs{Ack: ack}, &reply)
}

func (c *EngineRPCClient) EmitBatch(events []core.Event) error {
	var reply struct{}
	eventsJSON := make([][]byte, len(events))
	for index, event := range events {
		data, err := core.EncodeEvent(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event %d: %w", index, err)
		}
		eventsJSON[index] = data
	}
	return c.client.Call("Engine.EmitBatchRPC", &EmitBatchArgs{EventsJSON: eventsJSON}, &reply)
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
	return c.client.Call("Engine.SetStateRPC", &StateSetArgs{Key: key, Value: value}, &reply)
}

func (c *EngineRPCClient) SetStateWithTTL(key string, value []byte, ttl time.Duration) error {
	var reply struct{}
	return c.client.Call("Engine.SetStateWithTTLRPC", &StateSetWithTTLArgs{
		Key:      key,
		Value:    value,
		TTLNanos: ttl.Nanoseconds(),
	}, &reply)
}

func (c *EngineRPCClient) DeleteState(key string) error {
	var reply struct{}
	return c.client.Call("Engine.DeleteStateRPC", &StateKeyArgs{Key: key}, &reply)
}

func (c *EngineRPCClient) ExecuteConnector(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	var reply ConnectorExecuteReply
	if err := c.client.Call("Engine.ConnectorExecuteRPC", &req, &reply); err != nil {
		return core.ConnectorResponse{}, err
	}
	if reply.Error != "" {
		return core.ConnectorResponse{}, errors.New(reply.Error)
	}
	return reply.Response, nil
}

func (c *EngineRPCClient) CurrentResourceStatus(resourceRef string) (core.ResourceStatusEvent, bool) {
	var reply CurrentResourceStatusReply
	if err := c.client.Call("Engine.CurrentResourceStatusRPC", &CurrentResourceStatusArgs{
		ResourceRef: resourceRef,
	}, &reply); err != nil {
		return core.ResourceStatusEvent{}, false
	}
	return reply.Event, reply.Found
}

func (c *EngineRPCClient) NextResourceEvent(timeout time.Duration) (core.ResourceStatusEvent, bool, error) {
	var reply NextResourceEventReply
	if err := c.client.Call("Engine.NextResourceEventRPC", &NextResourceEventArgs{
		TimeoutMS: int(timeout / time.Millisecond),
	}, &reply); err != nil {
		return core.ResourceStatusEvent{}, false, err
	}
	return reply.Event, reply.OK, nil
}

var _ runtime.EngineBridge = (*EngineRPCClient)(nil)

func EdgeSubject(ruleID, fromNodeID, toNodeID string) string {
	return fmt.Sprintf("rule.edge.%s.%s.%s", ruleID, fromNodeID, toNodeID)
}
