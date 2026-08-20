package runtime

import "time"

type Telemetry struct {
	ExecutionID  string
	Runtime      string
	StartedAt    time.Time
	FinishedAt   time.Time
	Instructions uint64
	MemoryPeak   uint64
	OutputBytes  uint64
	CacheHit     bool
}

func (t Telemetry) Duration() time.Duration { return t.FinishedAt.Sub(t.StartedAt) }
func (t Telemetry) Valid() bool {
	return t.ExecutionID != "" && t.Runtime != "" && !t.StartedAt.IsZero() && t.FinishedAt.After(t.StartedAt)
}
func (t Telemetry) ValidStrict() bool { return t.Valid() && t.FinishedAt.After(t.StartedAt) }
