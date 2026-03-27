package sdk

import (
	"github.com/punk-one/punk-rule-plugin-sdk/internal/core"
	healthruntime "github.com/punk-one/punk-rule-plugin-sdk/internal/healthruntime"
)

type HealthStatus = core.HealthStatus

const (
	HealthUnknown      = core.HealthUnknown
	HealthInitializing = core.HealthInitializing
	HealthHealthy      = core.HealthHealthy
	HealthDegraded     = core.HealthDegraded
	HealthRecovering   = core.HealthRecovering
	HealthUnhealthy    = core.HealthUnhealthy
	HealthStopped      = core.HealthStopped
)

type HealthKind = core.HealthKind

const (
	HealthKindState     = core.HealthKindState
	HealthKindHeartbeat = core.HealthKindHeartbeat
)

type HealthOptions = core.HealthOptions

func DefaultHealthOptions() HealthOptions {
	return core.DefaultHealthOptions()
}

type HealthReport = core.HealthReport
type HealthHeartbeat = core.HealthHeartbeat
type ReportHealthArgs = core.ReportHealthArgs
type HealthReporter = core.HealthReporter

func newNoopHealthReporter() HealthReporter {
	return healthruntime.NewNoop()
}
