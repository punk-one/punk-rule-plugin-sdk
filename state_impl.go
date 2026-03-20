package sdk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/punk-rule-engine/punk-rule-engine/pkg/state"
)

// StateManagerImpl 生产级 StateManager 实现
type StateManagerImpl struct {
	ruleID  string
	nodeID  string
	manager *state.StateManager
}

// NewStateManager 创建 StateManager
func NewStateManager(ruleID, nodeID string, manager *state.StateManager) StateManager {
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
