package transport

import "github.com/hashicorp/go-plugin"

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PUNK_PLUGIN_MAGIC",
	MagicCookieValue: "punk-rule-engine",
}
