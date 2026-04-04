package runtime

import (
	"errors"
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
	ExecuteConnector(req core.ConnectorRequest) (core.ConnectorResponse, error)
	CurrentResourceStatus(resourceRef string) (core.ResourceStatusEvent, bool)
	NextResourceEvent(timeout time.Duration) (core.ResourceStatusEvent, bool, error)
}

type PluginRuntimeContext struct {
	ruleID         string
	nodeID         string
	engine         EngineBridge
	logger         core.Logger
	emitter        core.Emitter
	metrics        core.Metrics
	health         core.HealthReporter
	stateMgr       core.StateManager
	connector      core.ConnectorClient
	resourceEvents chan core.ResourceStatusEvent
	stopCh         chan struct{}
}

type noopConnectorClient struct{}
type engineConnectorClient struct {
	engine EngineBridge
}

func NewNoopConnectorClient() core.ConnectorClient {
	return noopConnectorClient{}
}

func (noopConnectorClient) Read(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	return core.ConnectorResponse{}, nil
}

func (noopConnectorClient) Write(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	return core.ConnectorResponse{}, nil
}

func (noopConnectorClient) BatchRead(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	return core.ConnectorResponse{}, nil
}

func (noopConnectorClient) BatchWrite(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	return core.ConnectorResponse{}, nil
}

func (noopConnectorClient) CurrentStatus(resourceRef string) (core.ResourceStatusEvent, bool) {
	return core.ResourceStatusEvent{}, false
}

func (c engineConnectorClient) Read(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	req.Kind = core.ConnectorRequestRead
	return c.execute(req)
}

func (c engineConnectorClient) Write(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	req.Kind = core.ConnectorRequestWrite
	return c.execute(req)
}

func (c engineConnectorClient) BatchRead(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	req.Kind = core.ConnectorRequestBatchRead
	return c.execute(req)
}

func (c engineConnectorClient) BatchWrite(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	req.Kind = core.ConnectorRequestBatchWrite
	return c.execute(req)
}

func (c engineConnectorClient) CurrentStatus(resourceRef string) (core.ResourceStatusEvent, bool) {
	if c.engine == nil {
		return core.ResourceStatusEvent{}, false
	}
	return c.engine.CurrentResourceStatus(resourceRef)
}

func (c engineConnectorClient) execute(req core.ConnectorRequest) (core.ConnectorResponse, error) {
	if c.engine == nil {
		return core.ConnectorResponse{}, errors.New("engine connector bridge is not available")
	}
	return c.engine.ExecuteConnector(req)
}

func NewPluginRuntimeContext(engine EngineBridge, ruleID, nodeID string, healthOptions core.HealthOptions, stateMgr core.StateManager) core.RuntimeContext {
	if stateMgr == nil {
		stateMgr = NewStateManager(ruleID, nodeID, NewRPCStateStore(engine))
	}
	ctx := &PluginRuntimeContext{
		ruleID:         ruleID,
		nodeID:         nodeID,
		engine:         engine,
		logger:         NewRPCLoggerFromEngine(engine, ruleID, nodeID),
		emitter:        NewRPCEmitterFromEngine(engine),
		metrics:        NewRPCMetricsFromEngine(engine, ruleID, nodeID),
		health:         NewEngineHealthReporter(engine, healthOptions),
		stateMgr:       stateMgr,
		connector:      engineConnectorClient{engine: engine},
		resourceEvents: make(chan core.ResourceStatusEvent, 8),
		stopCh:         make(chan struct{}),
	}
	if engine == nil {
		ctx.connector = NewNoopConnectorClient()
	} else {
		go ctx.resourceEventLoop()
	}
	return ctx
}

func (r *PluginRuntimeContext) RuleID() string                  { return r.ruleID }
func (r *PluginRuntimeContext) NodeID() string                  { return r.nodeID }
func (r *PluginRuntimeContext) Logger() core.Logger             { return r.logger }
func (r *PluginRuntimeContext) Emitter() core.Emitter           { return r.emitter }
func (r *PluginRuntimeContext) Metrics() core.Metrics           { return r.metrics }
func (r *PluginRuntimeContext) Health() core.HealthReporter     { return r.health }
func (r *PluginRuntimeContext) State() core.StateManager        { return r.stateMgr }
func (r *PluginRuntimeContext) Connector() core.ConnectorClient { return r.connector }
func (r *PluginRuntimeContext) ResourceEvents() <-chan core.ResourceStatusEvent {
	return r.resourceEvents
}

func (r *PluginRuntimeContext) Close() error {
	if r == nil || r.stopCh == nil {
		return nil
	}
	select {
	case <-r.stopCh:
	default:
		close(r.stopCh)
	}
	return nil
}

func (r *PluginRuntimeContext) resourceEventLoop() {
	if r == nil || r.engine == nil {
		return
	}
	for {
		select {
		case <-r.stopCh:
			close(r.resourceEvents)
			return
		default:
		}

		event, ok, err := r.engine.NextResourceEvent(time.Second)
		if err != nil || !ok {
			continue
		}

		select {
		case <-r.stopCh:
			close(r.resourceEvents)
			return
		case r.resourceEvents <- event:
		default:
			select {
			case <-r.resourceEvents:
			default:
			}
			r.resourceEvents <- event
		}
	}
}
