package sdk

import (
	"fmt"
	"sync"
	"time"
)

// PluginRuntimeContext 插件进程内的 RuntimeContext 实现
type PluginRuntimeContext struct {
	ruleID  string
	nodeID  string
	engine  EngineRPC
	logger  Logger
	emitter Emitter
	metrics Metrics
}

func NewPluginRuntimeContext(engine EngineRPC, ruleID, nodeID string) RuntimeContext {
	return &PluginRuntimeContext{
		ruleID:  ruleID,
		nodeID:  nodeID,
		engine:  engine,
		logger:  NewRPCLoggerFromEngine(engine, ruleID, nodeID),
		emitter: NewRPCEmitterFromEngine(engine),
		metrics: NewRPCMetricsFromEngine(engine, ruleID, nodeID),
	}
}

func (r *PluginRuntimeContext) RuleID() string   { return r.ruleID }
func (r *PluginRuntimeContext) NodeID() string   { return r.nodeID }
func (r *PluginRuntimeContext) Logger() Logger   { return r.logger }
func (r *PluginRuntimeContext) Emitter() Emitter { return r.emitter }
func (r *PluginRuntimeContext) Metrics() Metrics { return r.metrics }

// EngineRPCEmitter 基于 EngineRPC 的 Emitter
type EngineRPCEmitter struct {
	engine EngineRPC
}

func NewRPCEmitterFromEngine(engine EngineRPC) Emitter {
	return &EngineRPCEmitter{engine: engine}
}

func (e *EngineRPCEmitter) Publish(evt Event) error {
	return e.engine.Emit(evt)
}

func (e *EngineRPCEmitter) Ack(eventID string) error {
	return e.engine.Ack(eventID)
}

func (e *EngineRPCEmitter) EmitTo(label string, evt Event) error {
	if evt.Metadata == nil {
		evt.Metadata = make(map[string]string)
	}
	evt.Metadata["route_label"] = label
	return e.engine.Emit(evt)
}

func (e *EngineRPCEmitter) EmitBatch(events []Event) error {
	return e.engine.EmitBatch(events)
}

// EngineRPCLogger 基于 EngineRPC 的带缓冲 Logger
type EngineRPCLogger struct {
	engine   EngineRPC
	ruleID   string
	nodeID   string
	mu       sync.Mutex
	logCache map[LogLevel]*logBuffer
}

type logBuffer struct {
	messages []string
	fields   []map[string]interface{}
}

func NewRPCLoggerFromEngine(engine EngineRPC, ruleID, nodeID string) Logger {
	l := &EngineRPCLogger{
		engine:   engine,
		ruleID:   ruleID,
		nodeID:   nodeID,
		logCache: make(map[LogLevel]*logBuffer),
	}
	levels := []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}
	for _, lv := range levels {
		l.logCache[lv] = &logBuffer{
			messages: make([]string, 0, 100),
			fields:   make([]map[string]interface{}, 0, 100),
		}
	}
	go l.startFlusher()
	return l
}

func (l *EngineRPCLogger) startFlusher() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		l.Flush()
	}
}

func (l *EngineRPCLogger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for lv, buf := range l.logCache {
		if len(buf.messages) > 0 {
			l.engine.LogBatch(lv, buf.messages, buf.fields)
			buf.messages = buf.messages[:0]
			buf.fields = buf.fields[:0]
		}
	}
}

func (l *EngineRPCLogger) addLog(level LogLevel, msg string, fields ...Field) {
	fieldsMap := make(map[string]interface{})
	for _, f := range fields {
		fieldsMap[f.Key] = f.Value
	}
	l.mu.Lock()
	buf := l.logCache[level]
	buf.messages = append(buf.messages, msg)
	buf.fields = append(buf.fields, fieldsMap)
	if len(buf.messages) >= 100 {
		l.mu.Unlock()
		l.Flush()
		return
	}
	l.mu.Unlock()
}

func (l *EngineRPCLogger) Debug(msg string, fields ...Field) { l.addLog(LogLevelDebug, msg, fields...) }
func (l *EngineRPCLogger) Info(msg string, fields ...Field)  { l.addLog(LogLevelInfo, msg, fields...) }
func (l *EngineRPCLogger) Warn(msg string, fields ...Field)  { l.addLog(LogLevelWarn, msg, fields...) }
func (l *EngineRPCLogger) Error(msg string, fields ...Field) { l.addLog(LogLevelError, msg, fields...) }
func (l *EngineRPCLogger) LogBatch(level LogLevel, messages []string, fields []map[string]interface{}) {
	l.engine.LogBatch(level, messages, fields)
}

func NewRPCMetricsFromEngine(engine EngineRPC, ruleID, nodeID string) Metrics {
	return &DefaultMetrics{}
}

type NoOpEmitter struct{}

func (e *NoOpEmitter) Publish(evt Event) error { return fmt.Errorf("emitter not available") }
func (e *NoOpEmitter) EmitTo(label string, evt Event) error {
	return fmt.Errorf("emitter not available")
}
func (e *NoOpEmitter) Ack(eventID string) error       { return nil }
func (e *NoOpEmitter) EmitBatch(events []Event) error { return fmt.Errorf("emitter not available") }
