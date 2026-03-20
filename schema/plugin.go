package schema

// PluginDescriptor 定义插件的元数据
type PluginDescriptor struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // sink, source, function
	Role        string      `json:"role"` // "source", "processor", "sink", "control" (New spec)
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Inputs      []InputDesc `json:"inputs"` // 配置项定义
}

// InputDesc 定义插件配置项的元数据
type InputDesc struct {
	Key          string      `json:"key"`
	Type         string      `json:"type"`              // text, number, boolean, select, json
	Label        string      `json:"label"`             // 显示标签
	DefaultValue interface{} `json:"default"`           // 默认值（支持多种类型）
	Required     bool        `json:"required"`          // 是否必填
	Options      []string    `json:"options,omitempty"` // 用于 select 类型的选项
	Hint         string      `json:"hint,omitempty"`    // UI 上的 tooltip 提示
}

// Context 提供执行上下文（用于插件接口）
type Context struct {
	RuleID   string
	NodeID   string
	Config   map[string]interface{}
	IsDryRun bool
	IsDebug  bool
}
