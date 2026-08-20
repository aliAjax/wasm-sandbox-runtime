package observability

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	Requests       atomic.Uint64
	Errors         atomic.Uint64
	Executions     atomic.Uint64
	Succeeded      atomic.Uint64
	Failed         atomic.Uint64
	TimedOut       atomic.Uint64
	QueueWaitNanos atomic.Int64
}

func (m *Metrics) Request(ok bool) {
	m.Requests.Add(1)
	if !ok {
		m.Errors.Add(1)
	}
}
func (m *Metrics) Execution(state string) {
	m.Executions.Add(1)
	switch state {
	case "succeeded":
		m.Succeeded.Add(1)
	case "failed":
		m.Failed.Add(1)
	case "timed_out":
		m.TimedOut.Add(1)
	}
}
func (m *Metrics) QueueWait(d time.Duration) { m.QueueWaitNanos.Add(d.Nanoseconds()) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"requests": m.Requests.Load(), "errors": m.Errors.Load(), "executions": m.Executions.Load(), "succeeded": m.Succeeded.Load(), "failed": m.Failed.Load(), "timed_out": m.TimedOut.Load()}
}
