package runtime

import (
	"time"

	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
)

type EngineBridge interface {
	Emit(e core.Event) error
	EmitWithTargets(e core.Event, toNodeIDs []string) error
	Log(level core.LogLevel, msg string, fields map[string]interface{})
	LogBatch(level core.LogLevel, messages []string, fields []map[string]interface{})
	IncCounter(name string, labels map[string]string)
	Observe(name string, value float64, labels map[string]string)
	ReportHealth(args core.ReportHealthArgs) error
	Ack(eventID string) error
	PublishAck(ack core.AckMessage) error
	EmitBatch(events []core.Event) error
	GetState(key string) ([]byte, error)
	SetState(key string, value []byte) error
	SetStateWithTTL(key string, value []byte, ttl time.Duration) error
	DeleteState(key string) error
}

type PluginRuntimeContext struct {
	ruleID   string
	nodeID   string
	engine   EngineBridge
	logger   core.Logger
	emitter  core.Emitter
	metrics  core.Metrics
	health   core.HealthReporter
	stateMgr core.StateManager
}

func NewPluginRuntimeContext(engine EngineBridge, ruleID, nodeID string, healthOptions core.HealthOptions, stateMgr core.StateManager) core.RuntimeContext {
	if stateMgr == nil {
		stateMgr = NewStateManager(ruleID, nodeID, NewRPCStateStore(engine))
	}
	return &PluginRuntimeContext{
		ruleID:   ruleID,
		nodeID:   nodeID,
		engine:   engine,
		logger:   NewRPCLoggerFromEngine(engine, ruleID, nodeID),
		emitter:  NewRPCEmitterFromEngine(engine),
		metrics:  NewRPCMetricsFromEngine(engine, ruleID, nodeID),
		health:   NewEngineHealthReporter(engine, healthOptions),
		stateMgr: stateMgr,
	}
}

func (r *PluginRuntimeContext) RuleID() string              { return r.ruleID }
func (r *PluginRuntimeContext) NodeID() string              { return r.nodeID }
func (r *PluginRuntimeContext) Logger() core.Logger         { return r.logger }
func (r *PluginRuntimeContext) Emitter() core.Emitter       { return r.emitter }
func (r *PluginRuntimeContext) Metrics() core.Metrics       { return r.metrics }
func (r *PluginRuntimeContext) Health() core.HealthReporter { return r.health }
func (r *PluginRuntimeContext) State() core.StateManager    { return r.stateMgr }
