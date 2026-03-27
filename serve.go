package sdk

import transport "github.com/punk-one/punk-rule-plugin-sdk/internal/transport"

// ServeOptions 控制 sdk.Serve 的运行行为。
type ServeOptions struct {
	Health HealthOptions
}

// Serve 以 SDK 约定的默认配置启动插件进程。
func Serve(impl Plugin, options ServeOptions) {
	transport.Serve(impl, options.Health)
}
