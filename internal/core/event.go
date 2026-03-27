package core

import (
	"encoding/json"
	"fmt"
	"time"
)

const EventMetadataTraceIDKey = "trace_id"

type Event struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
	Metadata  map[string]string      `json:"metadata"`

	RuleID string `json:"-"`
	Source string `json:"-"`
}

func EncodeEvent(e Event) ([]byte, error) {
	return json.Marshal(e)
}

func DecodeEvent(data []byte) (Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	return e, err
}

func NewEvent(payload map[string]interface{}, metadata map[string]string) Event {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return Event{
		ID:        fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Payload:   payload,
		Metadata:  metadata,
	}
}

func (e Event) Clone() Event {
	clone := e
	clone.Payload = clonePayload(e.Payload)
	clone.Metadata = cloneMetadata(e.Metadata)
	return clone
}

func GetTraceID(e Event) string {
	if e.Metadata == nil {
		return ""
	}
	return e.Metadata[EventMetadataTraceIDKey]
}

func SetTraceID(e *Event, traceID string) {
	if e == nil {
		return
	}
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[EventMetadataTraceIDKey] = traceID
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return nil
	}
	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}

func clonePayload(payload map[string]interface{}) map[string]interface{} {
	if payload == nil {
		return nil
	}
	cloned := make(map[string]interface{}, len(payload))
	for key, value := range payload {
		cloned[key] = cloneValue(value)
	}
	return cloned
}

func cloneValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return clonePayload(typed)
	case []interface{}:
		cloned := make([]interface{}, len(typed))
		for index, item := range typed {
			cloned[index] = cloneValue(item)
		}
		return cloned
	default:
		return typed
	}
}
