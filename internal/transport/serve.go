package transport

import (
	"github.com/hashicorp/go-plugin"
	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
)

func Serve(impl core.Plugin, options core.HealthOptions) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"plugin": &PluginRPC{
				Impl:          impl,
				defaultHealth: options,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func ServeConnector(impl core.ConnectorPlugin) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"connector": &ConnectorRPC{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
