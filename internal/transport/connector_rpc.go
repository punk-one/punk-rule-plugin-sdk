package transport

import (
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
)

type ConnectorRPC struct {
	Impl core.ConnectorPlugin
}

func (p *ConnectorRPC) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &ConnectorRPCServer{Impl: p.Impl}, nil
}

func (p *ConnectorRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ConnectorRPCClient{Client: c}, nil
}

type ConnectorRPCServer struct {
	Impl core.ConnectorPlugin
}

func (s *ConnectorRPCServer) Info(args struct{}, reply *InfoReply) error {
	reply.Info = s.Impl.Info()
	return nil
}

func (s *ConnectorRPCServer) CreateResource(args *ConnectorCreateResourceArgs, reply *ConnectorCreateResourceReply) error {
	handle, err := s.Impl.CreateResource(args.Resource)
	reply.Handle = handle
	reply.Error = core.EncodePluginError(err)
	return nil
}

func (s *ConnectorRPCServer) DestroyResource(args *ConnectorDestroyResourceArgs, reply *ConnectorDestroyResourceReply) error {
	reply.Error = core.EncodePluginError(s.Impl.DestroyResource(args.ProviderHandle))
	return nil
}

func (s *ConnectorRPCServer) Execute(args *ConnectorExecuteArgs, reply *ConnectorExecuteReply) error {
	resp, err := s.Impl.Execute(args.ProviderHandle, args.Request)
	reply.Response = resp
	reply.Error = core.EncodePluginError(err)
	return nil
}

func (s *ConnectorRPCServer) Probe(args *ConnectorProbeArgs, reply *ConnectorProbeReply) error {
	event, err := s.Impl.Probe(args.ProviderHandle, args.Request)
	reply.Event = event
	reply.Error = core.EncodePluginError(err)
	return nil
}

func (s *ConnectorRPCServer) Stop(args struct{}, reply *StopReply) error {
	reply.Error = core.EncodePluginError(s.Impl.Stop())
	return nil
}

type ConnectorRPCClient struct {
	Client *rpc.Client
}

func (c *ConnectorRPCClient) Info() core.PluginInfo {
	var reply InfoReply
	if err := c.Client.Call("Plugin.Info", struct{}{}, &reply); err != nil {
		return core.PluginInfo{}
	}
	return reply.Info
}

func (c *ConnectorRPCClient) CreateResource(resource core.ConnectorResource) (string, error) {
	var reply ConnectorCreateResourceReply
	if err := callWithTimeout(c.Client, "Plugin.CreateResource", &ConnectorCreateResourceArgs{Resource: resource}, &reply, defaultRPCTimeout); err != nil {
		return "", err
	}
	return reply.Handle, core.DecodePluginError(reply.Error)
}

func (c *ConnectorRPCClient) DestroyResource(providerHandle string) error {
	var reply ConnectorDestroyResourceReply
	if err := callWithTimeout(c.Client, "Plugin.DestroyResource", &ConnectorDestroyResourceArgs{ProviderHandle: providerHandle}, &reply, defaultRPCTimeout); err != nil {
		return err
	}
	return core.DecodePluginError(reply.Error)
}

func (c *ConnectorRPCClient) Execute(providerHandle string, req core.ConnectorRequest) (core.ConnectorResponse, error) {
	return c.ExecuteWithTimeout(providerHandle, req, defaultRPCTimeout)
}

func (c *ConnectorRPCClient) ExecuteWithTimeout(providerHandle string, req core.ConnectorRequest, timeout time.Duration) (core.ConnectorResponse, error) {
	var reply ConnectorExecuteReply
	if err := callWithTimeout(c.Client, "Plugin.Execute", &ConnectorExecuteArgs{ProviderHandle: providerHandle, Request: req}, &reply, timeout); err != nil {
		return core.ConnectorResponse{}, err
	}
	return reply.Response, core.DecodePluginError(reply.Error)
}

func (c *ConnectorRPCClient) Probe(providerHandle string, req core.ConnectorRequest) (core.ResourceStatusEvent, error) {
	var reply ConnectorProbeReply
	if err := callWithTimeout(c.Client, "Plugin.Probe", &ConnectorProbeArgs{ProviderHandle: providerHandle, Request: req}, &reply, defaultRPCTimeout); err != nil {
		return core.ResourceStatusEvent{}, err
	}
	return reply.Event, core.DecodePluginError(reply.Error)
}

func (c *ConnectorRPCClient) Stop() error {
	var reply StopReply
	if err := callWithTimeout(c.Client, "Plugin.Stop", struct{}{}, &reply, defaultRPCTimeout); err != nil {
		return err
	}
	return core.DecodePluginError(reply.Error)
}
