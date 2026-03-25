package sdk

import "testing"

type engineRPCStub struct {
	state       map[string][]byte
	counter     string
	observation string
}

func newEngineRPCStub() *engineRPCStub {
	return &engineRPCStub{state: make(map[string][]byte)}
}

func (s *engineRPCStub) Emit(e Event) error                                { return nil }
func (s *engineRPCStub) EmitWithTargets(e Event, toNodeIDs []string) error { return nil }
func (s *engineRPCStub) Log(level LogLevel, msg string, fields map[string]interface{}) {
}
func (s *engineRPCStub) LogBatch(level LogLevel, messages []string, fields []map[string]interface{}) {
}
func (s *engineRPCStub) IncCounter(name string, labels map[string]string) { s.counter = name }
func (s *engineRPCStub) Observe(name string, value float64, labels map[string]string) {
	s.observation = name
}
func (s *engineRPCStub) Ack(eventID string) error            { return nil }
func (s *engineRPCStub) PublishAck(ack AckMessage) error     { return nil }
func (s *engineRPCStub) EmitBatch(events []Event) error      { return nil }
func (s *engineRPCStub) GetState(key string) ([]byte, error) { return s.state[key], nil }
func (s *engineRPCStub) SetState(key string, value []byte) error {
	s.state[key] = append([]byte(nil), value...)
	return nil
}
func (s *engineRPCStub) DeleteState(key string) error { delete(s.state, key); return nil }

func TestPluginRuntimeContextProvidesStateAndMetrics(t *testing.T) {
	engine := newEngineRPCStub()
	ctx := NewPluginRuntimeContext(engine, "rule-test", "node-test")

	stateful, ok := ctx.(interface{ State() StateManager })
	if !ok {
		t.Fatalf("expected plugin runtime context to expose State()")
	}
	if stateful.State() == nil {
		t.Fatalf("expected non-nil state manager")
	}

	type snapshot struct {
		Value int `json:"value"`
	}

	key := stateful.State().GenerateKey([]string{"device-1"})
	if err := stateful.State().SetState(key, &snapshot{Value: 42}); err != nil {
		t.Fatalf("SetState failed: %v", err)
	}

	var got snapshot
	if err := stateful.State().GetState(key, &got); err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	if got.Value != 42 {
		t.Fatalf("expected value 42, got %d", got.Value)
	}

	ctx.Metrics().IncCounter("events_processed", map[string]string{"node": "node-test"})
	ctx.Metrics().Observe("latency_ms", 12, nil)
	if engine.counter != "events_processed" {
		t.Fatalf("expected metric counter to be forwarded, got %q", engine.counter)
	}
	if engine.observation != "latency_ms" {
		t.Fatalf("expected observation to be forwarded, got %q", engine.observation)
	}
}
