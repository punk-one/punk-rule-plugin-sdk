package sdk

import "github.com/hashicorp/go-plugin"

// HandshakeConfig 是 go-plugin 的握手配置
// 确保插件和引擎使用的协议版本一致
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PUNK_PLUGIN_MAGIC",
	MagicCookieValue: "punk-rule-engine",
}
