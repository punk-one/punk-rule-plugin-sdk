package runtime

import (
	"sync"
	"time"

	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
)

type EngineRPCLogger struct {
	engine   EngineBridge
	ruleID   string
	nodeID   string
	mu       sync.Mutex
	logCache map[core.LogLevel]*logBuffer
}

type logBuffer struct {
	messages []string
	fields   []map[string]interface{}
}

func NewRPCLoggerFromEngine(engine EngineBridge, ruleID, nodeID string) core.Logger {
	logger := &EngineRPCLogger{
		engine:   engine,
		ruleID:   ruleID,
		nodeID:   nodeID,
		logCache: make(map[core.LogLevel]*logBuffer),
	}
	levels := []core.LogLevel{core.LogLevelDebug, core.LogLevelInfo, core.LogLevelWarn, core.LogLevelError}
	for _, level := range levels {
		logger.logCache[level] = &logBuffer{
			messages: make([]string, 0, 100),
			fields:   make([]map[string]interface{}, 0, 100),
		}
	}
	go logger.startFlusher()
	return logger
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
	for level, buf := range l.logCache {
		if len(buf.messages) == 0 {
			continue
		}
		l.engine.LogBatch(level, buf.messages, buf.fields)
		buf.messages = buf.messages[:0]
		buf.fields = buf.fields[:0]
	}
}

func (l *EngineRPCLogger) addLog(level core.LogLevel, msg string, fields ...core.Field) {
	fieldsMap := make(map[string]interface{})
	for _, field := range fields {
		fieldsMap[field.Key] = field.Value
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

func (l *EngineRPCLogger) Debug(msg string, fields ...core.Field) {
	l.addLog(core.LogLevelDebug, msg, fields...)
}

func (l *EngineRPCLogger) Info(msg string, fields ...core.Field) {
	l.addLog(core.LogLevelInfo, msg, fields...)
}

func (l *EngineRPCLogger) Warn(msg string, fields ...core.Field) {
	l.addLog(core.LogLevelWarn, msg, fields...)
}

func (l *EngineRPCLogger) Error(msg string, fields ...core.Field) {
	l.addLog(core.LogLevelError, msg, fields...)
}

func (l *EngineRPCLogger) LogBatch(level core.LogLevel, messages []string, fields []map[string]interface{}) {
	l.engine.LogBatch(level, messages, fields)
}
