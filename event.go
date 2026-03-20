package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event 规则引擎统一事件定义 (v1.1.0)
type Event struct {
	ID        string                 `json:"id"`        // 事件唯一 ID
	Timestamp time.Time              `json:"timestamp"` // 事件产生时间
	Payload   map[string]interface{} `json:"payload"`   // 数据载荷
	Metadata  map[string]string      `json:"metadata"`  // 元数据 (trace_id, source等)

	// 以下字段仅内部流转使用，不序列化到外部
	RuleID string `json:"-"`
	Source string `json:"-"`
}

// EncodeEvent 编码事件为字节数组
func EncodeEvent(e Event) ([]byte, error) {
	return json.Marshal(e)
}

// DecodeEvent 从字节数组解码事件
func DecodeEvent(data []byte) (Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	return e, err
}

// NewEvent 创建新事件
// 自动生成 ID 和时间戳
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
