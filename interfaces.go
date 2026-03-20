package sdk

// Logger 日志接口（v1.1.1）
// 插件通过 Logger 记录日志，而不是直接使用标准库
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	// LogBatch 批量记录日志 (v1.5.0)
	LogBatch(level LogLevel, messages []string, fields []map[string]interface{})
}

// Emitter 事件发布能力接口（v1.1.1）
type Emitter interface {
	Publish(event Event) error
	Ack(eventID string) error
	// EmitTo 发送到指定路由标签（v1.3.1）
	EmitTo(label string, event Event) error
	// EmitBatch 批量发布事件 (v1.5.0)
	EmitBatch(events []Event) error
}

// Metrics 指标接口（v1.1.1）
type Metrics interface {
	// IncCounter 增加计数器
	IncCounter(name string, labels map[string]string)
	// Observe 记录观测值（如耗时、大小）
	Observe(name string, value float64, labels map[string]string)
}

// StateManager 状态管理接口（v1.1.1）
type StateManager interface {
	// GetState 获取状态
	GetState(key string, state interface{}) error
	// SetState 设置状态
	SetState(key string, state interface{}) error
	// DeleteState 删除状态
	DeleteState(key string) error
	// GenerateKey 生成符合规范的 State Key
	GenerateKey(keys []string) string
}

// RuntimeContext 插件运行上下文（v1.1.1）
type RuntimeContext interface {
	// RuleID 返回所属规则 ID
	RuleID() string
	// NodeID 返回当前节点 ID
	NodeID() string
	// Logger 返回日志记录器
	Logger() Logger
	// Emitter 返回事件发布器
	Emitter() Emitter
	// Metrics 返回指标记录器
	Metrics() Metrics
}

// Plugin 插件核心接口 (v1.1.0)
type Plugin interface {
	// Info 返回插件能力声明
	Info() PluginInfo

	// Init 初始化插件（配置校验）
	// 配置已通过 Schema 校验，插件只需解析和存储
	Init(cfg PluginConfig) error

	// Start 启动插件（开始处理数据）
	// Source Plugin: 调用 ctx.Emitter().Publish() 发布事件
	// Processor/Sink Plugin: 通过 OnEvent 接收事件（由 Dispatcher 调用）
	Start(ctx RuntimeContext) error

	// OnEvent 接收事件 (v1.2)
	// Dispatcher 通过此方法向插件推送事件
	OnEvent(e Event) error

	// OnEvents 批量接收事件 (v1.5.0)
	// 如果插件支持批量处理，应实现此方法以提高性能
	OnEvents(events []Event) error

	// Stop 停止插件（清理资源）
	Stop() error
}
