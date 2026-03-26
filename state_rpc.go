package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type rpcStateStore struct {
	engine EngineRPC
}

func newRPCStateStore(engine EngineRPC) StateStore {
	return &rpcStateStore{engine: engine}
}

func (s *rpcStateStore) GetState(ctx context.Context, key string, state interface{}) error {
	if s.engine == nil {
		return fmt.Errorf("state manager not available")
	}
	data, err := s.engine.GetState(key)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("state not found: %s", key)
	}
	return json.Unmarshal(data, state)
}

func (s *rpcStateStore) SetState(ctx context.Context, key string, state interface{}) error {
	return s.SetStateWithTTL(ctx, key, state, 0)
}

func (s *rpcStateStore) SetStateWithTTL(ctx context.Context, key string, state interface{}, ttl time.Duration) error {
	if s.engine == nil {
		return fmt.Errorf("state manager not available")
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return s.engine.SetStateWithTTL(key, data, ttl)
}

func (s *rpcStateStore) DeleteState(ctx context.Context, key string) error {
	if s.engine == nil {
		return fmt.Errorf("state manager not available")
	}
	return s.engine.DeleteState(key)
}
