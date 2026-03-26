package sdk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// StateManagerImpl 生产级 StateManager 实现
type StateManagerImpl struct {
	ruleID  string
	nodeID  string
	manager StateStore
}

// StateStore 定义状态存储后端最小能力，避免 SDK 反向依赖 engine 内部包。
type StateStore interface {
	GetState(ctx context.Context, key string, state interface{}) error
	SetState(ctx context.Context, key string, state interface{}) error
	SetStateWithTTL(ctx context.Context, key string, state interface{}, ttl time.Duration) error
	DeleteState(ctx context.Context, key string) error
}

// NewStateManager 创建 StateManager
func NewStateManager(ruleID, nodeID string, manager StateStore) StateManager {
	return &StateManagerImpl{
		ruleID:  ruleID,
		nodeID:  nodeID,
		manager: manager,
	}
}

func (s *StateManagerImpl) GetState(key string, state interface{}) error {
	return s.manager.GetState(context.Background(), key, state)
}

func (s *StateManagerImpl) SetState(key string, state interface{}) error {
	return s.manager.SetState(context.Background(), key, state)
}

func (s *StateManagerImpl) SetStateWithTTL(key string, state interface{}, ttl time.Duration) error {
	return s.manager.SetStateWithTTL(context.Background(), key, state, ttl)
}

func (s *StateManagerImpl) DeleteState(key string) error {
	return s.manager.DeleteState(context.Background(), key)
}

// GenerateKey 生成 State Key（规范 3.2）
func (s *StateManagerImpl) GenerateKey(keys []string) string {
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	// 拼接：rule_id:node_id:key1:key2:...
	keyStr := fmt.Sprintf("%s:%s:%s", s.ruleID, s.nodeID, strings.Join(sortedKeys, ":"))

	// 计算 hash
	hash := sha256.Sum256([]byte(keyStr))
	return fmt.Sprintf("state:%s", hex.EncodeToString(hash[:]))
}
