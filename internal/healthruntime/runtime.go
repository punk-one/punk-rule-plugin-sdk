package healthruntime

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"
)

type Status string

const (
	StatusUnknown      Status = "unknown"
	StatusInitializing Status = "initializing"
	StatusHealthy      Status = "healthy"
	StatusDegraded     Status = "degraded"
	StatusRecovering   Status = "recovering"
	StatusUnhealthy    Status = "unhealthy"
	StatusStopped      Status = "stopped"
)

type Kind string

const (
	KindState     Kind = "state"
	KindHeartbeat Kind = "heartbeat"
)

type Options struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval,omitempty"`
	MaxSilencePeriod  time.Duration `json:"max_silence_period,omitempty"`
	QueueCapacity     int           `json:"queue_capacity,omitempty"`
}

func DefaultOptions() Options {
	return Options{
		HeartbeatInterval: 15 * time.Second,
		MaxSilencePeriod:  time.Minute,
		QueueCapacity:     32,
	}
}

func (o Options) Normalize() Options {
	defaults := DefaultOptions()
	if o.HeartbeatInterval <= 0 {
		o.HeartbeatInterval = defaults.HeartbeatInterval
	}
	if o.MaxSilencePeriod <= 0 {
		o.MaxSilencePeriod = defaults.MaxSilencePeriod
	}
	if o.QueueCapacity <= 0 {
		o.QueueCapacity = defaults.QueueCapacity
	}
	return o
}

type Report struct {
	Status     Status            `json:"status"`
	Code       string            `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	Details    map[string]string `json:"details,omitempty"`
	ObservedAt time.Time         `json:"observed_at,omitempty"`
	Sequence   uint64            `json:"sequence,omitempty"`
}

type Heartbeat struct {
	Message    string            `json:"message,omitempty"`
	Details    map[string]string `json:"details,omitempty"`
	ObservedAt time.Time         `json:"observed_at,omitempty"`
	Sequence   uint64            `json:"sequence,omitempty"`
}

type ReportArgs struct {
	Kind               Kind              `json:"kind"`
	Status             Status            `json:"status,omitempty"`
	Code               string            `json:"code,omitempty"`
	Message            string            `json:"message,omitempty"`
	Details            map[string]string `json:"details,omitempty"`
	ObservedAtUnixNano int64             `json:"observed_at_unix_nano,omitempty"`
	Sequence           uint64            `json:"sequence,omitempty"`
}

type Reporter interface {
	Report(report Report) error
	Healthy(msg string, details map[string]string) error
	Degraded(code, msg string, details map[string]string) error
	Recovering(code, msg string, details map[string]string) error
	Unhealthy(code, msg string, details map[string]string) error
	Heartbeat(msg string, details map[string]string) error
}

type Sender interface {
	SendHealth(args ReportArgs) error
}

type healthEnvelope struct {
	args ReportArgs
}

type asyncReporter struct {
	sender  Sender
	options Options

	queue  chan healthEnvelope
	stopCh chan struct{}

	sequence atomic.Uint64

	mu            sync.Mutex
	lastStateKey  string
	lastStateAt   time.Time
	lastHeartbeat time.Time
	closed        bool
}

type nopReporter struct{}

func (nopReporter) Report(Report) error                                { return nil }
func (nopReporter) Healthy(string, map[string]string) error            { return nil }
func (nopReporter) Degraded(string, string, map[string]string) error   { return nil }
func (nopReporter) Recovering(string, string, map[string]string) error { return nil }
func (nopReporter) Unhealthy(string, string, map[string]string) error  { return nil }
func (nopReporter) Heartbeat(string, map[string]string) error          { return nil }

func NewNoop() Reporter {
	return nopReporter{}
}

func NewAsync(sender Sender, options Options) Reporter {
	if sender == nil {
		return NewNoop()
	}

	reporter := &asyncReporter{
		sender:  sender,
		options: options.Normalize(),
		stopCh:  make(chan struct{}),
	}
	reporter.queue = make(chan healthEnvelope, reporter.options.QueueCapacity)
	go reporter.run()
	return reporter
}

func (r *asyncReporter) Close() error {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil
	}
	r.closed = true
	close(r.stopCh)
	r.mu.Unlock()
	return nil
}

func (r *asyncReporter) Report(report Report) error {
	if report.Status == "" {
		report.Status = StatusUnknown
	}
	if report.ObservedAt.IsZero() {
		report.ObservedAt = time.Now()
	}
	report.Details = cloneDetails(report.Details)
	report.Sequence = r.sequence.Add(1)

	stateKey := buildStateKey(report)

	r.mu.Lock()
	if !r.lastStateAt.IsZero() &&
		r.lastStateKey == stateKey &&
		report.ObservedAt.Sub(r.lastStateAt) < r.options.MaxSilencePeriod {
		r.mu.Unlock()
		return nil
	}
	r.lastStateKey = stateKey
	r.lastStateAt = report.ObservedAt
	r.lastHeartbeat = report.ObservedAt
	r.mu.Unlock()

	return r.enqueue(ReportArgs{
		Kind:               KindState,
		Status:             report.Status,
		Code:               report.Code,
		Message:            report.Message,
		Details:            report.Details,
		ObservedAtUnixNano: report.ObservedAt.UnixNano(),
		Sequence:           report.Sequence,
	})
}

func (r *asyncReporter) Healthy(msg string, details map[string]string) error {
	return r.Report(Report{
		Status:  StatusHealthy,
		Message: msg,
		Details: details,
	})
}

func (r *asyncReporter) Degraded(code, msg string, details map[string]string) error {
	return r.Report(Report{
		Status:  StatusDegraded,
		Code:    code,
		Message: msg,
		Details: details,
	})
}

func (r *asyncReporter) Recovering(code, msg string, details map[string]string) error {
	return r.Report(Report{
		Status:  StatusRecovering,
		Code:    code,
		Message: msg,
		Details: details,
	})
}

func (r *asyncReporter) Unhealthy(code, msg string, details map[string]string) error {
	return r.Report(Report{
		Status:  StatusUnhealthy,
		Code:    code,
		Message: msg,
		Details: details,
	})
}

func (r *asyncReporter) Heartbeat(msg string, details map[string]string) error {
	observedAt := time.Now()
	details = cloneDetails(details)

	r.mu.Lock()
	if !r.lastHeartbeat.IsZero() && observedAt.Sub(r.lastHeartbeat) < r.options.HeartbeatInterval {
		r.mu.Unlock()
		return nil
	}
	r.lastHeartbeat = observedAt
	r.mu.Unlock()

	return r.enqueue(ReportArgs{
		Kind:               KindHeartbeat,
		Message:            msg,
		Details:            details,
		ObservedAtUnixNano: observedAt.UnixNano(),
		Sequence:           r.sequence.Add(1),
	})
}

func (r *asyncReporter) enqueue(args ReportArgs) error {
	select {
	case <-r.stopCh:
		return nil
	default:
	}

	select {
	case r.queue <- healthEnvelope{args: args}:
	default:
	}
	return nil
}

func (r *asyncReporter) run() {
	ticker := time.NewTicker(r.options.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopCh:
			return
		case item := <-r.queue:
			_ = r.sender.SendHealth(item.args)
		case <-ticker.C:
			_ = r.Heartbeat("", nil)
		}
	}
}

func buildStateKey(report Report) string {
	detailsJSON, _ := json.Marshal(report.Details)
	return string(report.Status) + "|" + report.Code + "|" + report.Message + "|" + string(detailsJSON)
}

func cloneDetails(details map[string]string) map[string]string {
	if len(details) == 0 {
		return nil
	}
	out := make(map[string]string, len(details))
	for key, value := range details {
		out[key] = value
	}
	return out
}
