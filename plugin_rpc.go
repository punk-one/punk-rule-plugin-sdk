package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
)

// PluginRPC 实现了 go-plugin 的 Plugin 接口
type PluginRPC struct {
	Impl Plugin
}

func (p *PluginRPC) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &PluginRPCServer{Impl: p.Impl, broker: b}, nil
}

func (p *PluginRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginRPCClient{Client: c, Broker: b}, nil
}

// PluginRPCServer RPC 服务端
type PluginRPCServer struct {
	Impl   Plugin
	broker *plugin.MuxBroker
}

func (s *PluginRPCServer) Info(args struct{}, reply *InfoReply) error {
	reply.Info = s.Impl.Info()
	return nil
}

func (s *PluginRPCServer) Init(args *InitArgs, reply *InitReply) error {
	reply.Error = errToString(s.Impl.Init(args.Config))
	return nil
}

func (s *PluginRPCServer) Start(args *StartArgs, reply *StartReply) error {
	// 使用 MuxBroker 建立 EngineRPC 连接（ID: 3）
	conn, err := s.broker.Dial(3)
	if err != nil {
		reply.Error = fmt.Sprintf("failed to dial Engine RPC connection: %v", err)
		return nil
	}

	// 创建 RPC client
	rpcClient := rpc.NewClient(conn)

	// 创建 EngineRPC client
	engineRPC := NewEngineRPCClient(rpcClient)

	// 创建 RuntimeContext
	ctx := NewPluginRuntimeContext(engineRPC, args.RuleID, args.NodeID)

	// 调用插件实现的 Start 方法
	if err := s.Impl.Start(ctx); err != nil {
		reply.Error = err.Error()
		return nil
	}

	return nil
}

func (s *PluginRPCServer) ReceiveEvent(args *ReceiveEventArgs, reply *ReceiveEventReply) error {
	var event Event
	if err := json.Unmarshal(args.EventJSON, &event); err != nil {
		reply.Error = fmt.Sprintf("failed to unmarshal event: %v", err)
		return nil
	}
	reply.Error = errToString(s.Impl.OnEvent(event))
	return nil
}

func (s *PluginRPCServer) ReceiveEvents(args *ReceiveEventsArgs, reply *ReceiveEventsReply) error {
	var events []Event
	if err := json.Unmarshal(args.EventsJSON, &events); err != nil {
		reply.Error = fmt.Sprintf("failed to unmarshal events: %v", err)
		return nil
	}
	reply.Error = errToString(s.Impl.OnEvents(events))
	return nil
}

func (s *PluginRPCServer) Stop(args struct{}, reply *StopReply) error {
	reply.Error = errToString(s.Impl.Stop())
	return nil
}

// PluginRPCClient RPC 客户端
type PluginRPCClient struct {
	Client *rpc.Client
	Broker *plugin.MuxBroker
}

const defaultRPCTimeout = 30 * time.Second

func (c *PluginRPCClient) Info() PluginInfo {
	var reply InfoReply
	err := c.Client.Call("Plugin.Info", struct{}{}, &reply)
	if err != nil {
		return PluginInfo{}
	}
	return reply.Info
}

func (c *PluginRPCClient) Init(cfg PluginConfig) error {
	return c.InitWithTimeout(cfg, defaultRPCTimeout)
}

func (c *PluginRPCClient) InitWithTimeout(cfg PluginConfig, timeout time.Duration) error {
	args := &InitArgs{Config: cfg}
	var reply InitReply
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- c.Client.Call("Plugin.Init", args, &reply) }()
	select {
	case err := <-done:
		if err != nil {
			return err
		}
		if reply.Error != "" {
			return fmt.Errorf("%s", reply.Error)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("RPC timeout: Init")
	}
}

func (c *PluginRPCClient) Start(args *StartArgs) error {
	var reply StartReply
	return c.Client.Call("Plugin.Start", args, &reply)
}

func (c *PluginRPCClient) ReceiveEvent(e Event) error {
	return c.ReceiveEventWithTimeout(e, defaultRPCTimeout)
}

func (c *PluginRPCClient) ReceiveEventWithTimeout(e Event, timeout time.Duration) error {
	eventJSON, err := json.Marshal(e)
	if err != nil {
		return err
	}
	args := &ReceiveEventArgs{EventJSON: eventJSON}
	var reply ReceiveEventReply
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- c.Client.Call("Plugin.ReceiveEvent", args, &reply) }()
	select {
	case err := <-done:
		if err != nil {
			return err
		}
		if reply.Error != "" {
			return fmt.Errorf("%s", reply.Error)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("RPC timeout: ReceiveEvent")
	}
}

func (c *PluginRPCClient) ReceiveEvents(events []Event) error {
	return c.ReceiveEventsWithTimeout(events, defaultRPCTimeout)
}

func (c *PluginRPCClient) ReceiveEventsWithTimeout(events []Event, timeout time.Duration) error {
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		return err
	}
	args := &ReceiveEventsArgs{EventsJSON: eventsJSON}
	var reply ReceiveEventsReply
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- c.Client.Call("Plugin.ReceiveEvents", args, &reply) }()
	select {
	case err := <-done:
		if err != nil {
			return err
		}
		if reply.Error != "" {
			return fmt.Errorf("%s", reply.Error)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("RPC timeout: ReceiveEvents")
	}
}

func (c *PluginRPCClient) Stop() error {
	return c.StopWithTimeout(defaultRPCTimeout)
}

func (c *PluginRPCClient) StopWithTimeout(timeout time.Duration) error {
	var reply StopReply
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- c.Client.Call("Plugin.Stop", struct{}{}, &reply) }()
	select {
	case err := <-done:
		if err != nil {
			return err
		}
		if reply.Error != "" {
			return fmt.Errorf("%s", reply.Error)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("RPC timeout: Stop")
	}
}

// RPC Args and Replies

type InfoReply struct{ Info PluginInfo }
type InitArgs struct{ Config PluginConfig }
type InitReply struct{ Error string }
type StartArgs struct {
	RuleID string
	NodeID string
}
type StartReply struct{ Error string }
type ReceiveEventArgs struct{ EventJSON []byte }
type ReceiveEventReply struct{ Error string }
type ReceiveEventsArgs struct{ EventsJSON []byte }
type ReceiveEventsReply struct{ Error string }
type StopReply struct{ Error string }

func errToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
