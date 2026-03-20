package sdk

import (
	"context"
	"fmt"
	"log"
)

// DefaultLogger 默认日志实现
// 使用标准库 log 包，避免依赖外部包
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *DefaultLogger) Info(msg string, fields ...Field) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (l *DefaultLogger) Error(msg string, fields ...Field) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *DefaultLogger) LogBatch(level LogLevel, messages []string, fields []map[string]interface{}) {
	for i, msg := range messages {
		var f []Field
		if i < len(fields) {
			for k, v := range fields[i] {
				f = append(f, Field{Key: k, Value: v})
			}
		}
		switch level {
		case LogLevelDebug:
			l.Debug(msg, f...)
		case LogLevelInfo:
			l.Info(msg, f...)
		case LogLevelWarn:
			l.Warn(msg, f...)
		case LogLevelError:
			l.Error(msg, f...)
		}
	}
}

// DefaultMetrics 默认指标实现（空实现）
type DefaultMetrics struct{}

func (m *DefaultMetrics) IncCounter(name string, labels map[string]string)             {}
func (m *DefaultMetrics) Observe(name string, value float64, labels map[string]string) {}

// EmitterAdapter 将 Publish 函数适配为 Emitter 接口
type EmitterAdapter struct {
	publisher func(port int, event Event) error
}

func (e *EmitterAdapter) Publish(event Event) error {
	if e.publisher == nil {
		return fmt.Errorf("publisher not configured")
	}
	return e.publisher(0, event)
}

func (e *EmitterAdapter) Ack(eventID string) error {
	return nil
}

func (e *EmitterAdapter) EmitTo(label string, event Event) error {
	if e.publisher == nil {
		return fmt.Errorf("publisher not configured")
	}
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["route_label"] = label
	return e.publisher(0, event)
}

func (e *EmitterAdapter) EmitBatch(events []Event) error {
	if e.publisher == nil {
		return fmt.Errorf("publisher not configured")
	}
	for _, event := range events {
		if err := e.publisher(0, event); err != nil {
			return err
		}
	}
	return nil
}

// RuntimeContextImpl RuntimeContext 实现（用于内置插件）
type RuntimeContextImpl struct {
	ruleID          string
	nodeID          string
	ctx             context.Context
	logger          Logger
	metrics         Metrics
	emitter         Emitter
	stateMgr        StateManager
	saveEngineRPCCB func(engineRPC interface{})
}

func (r *RuntimeContextImpl) RuleID() string           { return r.ruleID }
func (r *RuntimeContextImpl) NodeID() string           { return r.nodeID }
func (r *RuntimeContextImpl) Logger() Logger           { return r.logger }
func (r *RuntimeContextImpl) Emitter() Emitter         { return r.emitter }
func (r *RuntimeContextImpl) Metrics() Metrics         { return r.metrics }
func (r *RuntimeContextImpl) Context() context.Context { return r.ctx }
func (r *RuntimeContextImpl) State() StateManager      { return r.stateMgr }

func (r *RuntimeContextImpl) SaveEngineRPC(engineRPC interface{}) {
	if r.saveEngineRPCCB != nil {
		r.saveEngineRPCCB(engineRPC)
	}
}
