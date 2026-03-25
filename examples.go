package sdk

import (
	"encoding/json"
)

// This file contains plugin development examples and best practices

// ExamplePlugin demonstrates how to implement a complete processor plugin
type ExamplePlugin struct {
	config map[string]interface{}
	ctx    RuntimeContext
}

func (p *ExamplePlugin) Info() PluginInfo {
	return PluginInfo{
		ID:          "example-processor",
		Name:        "Example Processor",
		Version:     "1.5.2",
		Description: "An example processor plugin",
		Type:        PluginTypeProcessor,
		Capabilities: PluginCapabilities{
			SupportBatch: true,
			Stateful:     false,
		},
	}
}

func (p *ExamplePlugin) Init(cfg PluginConfig) error {
	if err := json.Unmarshal(cfg.Raw, &p.config); err != nil {
		return err
	}
	return nil
}

func (p *ExamplePlugin) Start(ctx RuntimeContext) error {
	p.ctx = ctx
	return nil
}

func (p *ExamplePlugin) OnEvent(e Event) error {
	// Processing logic
	e.Payload["processed"] = true
	return p.ctx.Emitter().Publish(e)
}

func (p *ExamplePlugin) OnEvents(events []Event) error {
	for i := range events {
		events[i].Payload["processed"] = true
	}
	return p.ctx.Emitter().EmitBatch(events)
}

func (p *ExamplePlugin) Stop() error {
	return nil
}

// ExampleSourcePlugin demonstrates a source plugin
type ExampleSourcePlugin struct {
	config map[string]interface{}
	ctx    RuntimeContext
}

func (p *ExampleSourcePlugin) Info() PluginInfo {
	return PluginInfo{
		ID:      "example-source",
		Name:    "Example Source",
		Version: "1.5.2",
		Type:    PluginTypeSource,
	}
}

func (p *ExampleSourcePlugin) Init(cfg PluginConfig) error {
	if err := json.Unmarshal(cfg.Raw, &p.config); err != nil {
		return err
	}
	return nil
}

func (p *ExampleSourcePlugin) Start(ctx RuntimeContext) error {
	p.ctx = ctx
	go func() {
		// Mock data generation
		event := NewEvent(map[string]interface{}{"data": "value"}, nil)
		_ = p.ctx.Emitter().Publish(event)
	}()
	return nil
}

func (p *ExampleSourcePlugin) OnEvent(e Event) error         { return nil } // Source doesn't receive events
func (p *ExampleSourcePlugin) OnEvents(events []Event) error { return nil }

func (p *ExampleSourcePlugin) Stop() error { return nil }
