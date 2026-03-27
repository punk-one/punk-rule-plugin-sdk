package runtime

import (
	"testing"
	"time"

	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
)

type engineBridgeStub struct {
	state       map[string][]byte
	counter     string
	observation string
	health      []core.ReportHealthArgs
}

func newEngineBridgeStub() *engineBridgeStub {
	return &engineBridgeStub{state: make(map[string][]byte)}
}

func (s *engineBridgeStub) Emit(e core.Event) error                                { return nil }
func (s *engineBridgeStub) EmitWithTargets(e core.Event, toNodeIDs []string) error { return nil }
func (s *engineBridgeStub) Log(level core.LogLevel, msg string, fields map[string]interface{}) {
}
func (s *engineBridgeStub) LogBatch(level core.LogLevel, messages []string, fields []map[string]interface{}) {
}
func (s *engineBridgeStub) IncCounter(name string, labels map[string]string) { s.counter = name }
func (s *engineBridgeStub) Observe(name string, value float64, labels map[string]string) {
	s.observation = name
}
func (s *engineBridgeStub) ReportHealth(args core.ReportHealthArgs) error {
	s.health = append(s.health, args)
	return nil
}
func (s *engineBridgeStub) Ack(eventID string) error             { return nil }
func (s *engineBridgeStub) PublishAck(ack core.AckMessage) error { return nil }
func (s *engineBridgeStub) EmitBatch(events []core.Event) error  { return nil }
func (s *engineBridgeStub) GetState(key string) ([]byte, error)  { return s.state[key], nil }
func (s *engineBridgeStub) DeleteState(key string) error         { delete(s.state, key); return nil }
func (s *engineBridgeStub) SetStateWithTTL(key string, value []byte, ttl time.Duration) error {
	return s.SetState(key, value)
}
func (s *engineBridgeStub) SetState(key string, value []byte) error {
	s.state[key] = append([]byte(nil), value...)
	return nil
}

func waitForHealthMessages(t *testing.T, stub *engineBridgeStub, expected int) []core.ReportHealthArgs {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if len(stub.health) >= expected {
			out := make([]core.ReportHealthArgs, len(stub.health))
			copy(out, stub.health)
			return out
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("expected at least %d health messages, got %d", expected, len(stub.health))
	return nil
}

func TestNewPluginRuntimeContextProvidesStateAndMetrics(t *testing.T) {
	engine := newEngineBridgeStub()
	ctx := NewPluginRuntimeContext(engine, "rule-test", "node-test", core.DefaultHealthOptions(), nil)

	if ctx.RuleID() != "rule-test" || ctx.NodeID() != "node-test" {
		t.Fatalf("unexpected context identity: %s/%s", ctx.RuleID(), ctx.NodeID())
	}

	stateful, ok := ctx.(interface{ State() core.StateManager })
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

	if closer, ok := ctx.Health().(interface{ Close() error }); ok {
		t.Cleanup(func() { _ = closer.Close() })
	}
	if err := ctx.Health().Healthy("ok", map[string]string{"phase": "test"}); err != nil {
		t.Fatalf("expected health reporter to be available: %v", err)
	}
}

func TestHealthReporterSuppressesDuplicateState(t *testing.T) {
	engine := newEngineBridgeStub()
	ctx := NewPluginRuntimeContext(engine, "rule-test", "node-test", core.HealthOptions{
		HeartbeatInterval: time.Hour,
		MaxSilencePeriod:  time.Minute,
		QueueCapacity:     8,
	}, nil)
	if closer, ok := ctx.Health().(interface{ Close() error }); ok {
		t.Cleanup(func() { _ = closer.Close() })
	}

	if err := ctx.Health().Healthy("connected", map[string]string{"target": "plc-1"}); err != nil {
		t.Fatalf("first report failed: %v", err)
	}
	if err := ctx.Health().Healthy("connected", map[string]string{"target": "plc-1"}); err != nil {
		t.Fatalf("second report failed: %v", err)
	}

	got := waitForHealthMessages(t, engine, 1)
	if len(got) != 1 {
		t.Fatalf("expected 1 report after suppression, got %d", len(got))
	}
	if got[0].Kind != core.HealthKindState {
		t.Fatalf("expected state report, got %s", got[0].Kind)
	}
	if got[0].Status != core.HealthHealthy {
		t.Fatalf("expected healthy status, got %s", got[0].Status)
	}
}

func TestHealthReporterEmitsHeartbeat(t *testing.T) {
	engine := newEngineBridgeStub()
	ctx := NewPluginRuntimeContext(engine, "rule-test", "node-test", core.HealthOptions{
		HeartbeatInterval: 25 * time.Millisecond,
		MaxSilencePeriod:  time.Minute,
		QueueCapacity:     8,
	}, nil)
	if closer, ok := ctx.Health().(interface{ Close() error }); ok {
		t.Cleanup(func() { _ = closer.Close() })
	}

	got := waitForHealthMessages(t, engine, 1)
	if got[0].Kind != core.HealthKindHeartbeat {
		t.Fatalf("expected heartbeat report, got %s", got[0].Kind)
	}
}
