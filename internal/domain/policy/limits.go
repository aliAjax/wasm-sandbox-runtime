package policy

import "time"

func (l Limits) Timeout() time.Duration        { return time.Duration(l.TimeoutSeconds) * time.Second }
func (l Limits) RetryAllowed(attempt int) bool { return attempt >= 0 && attempt < l.Retries }
func (l Limits) ClampCPU(v uint64) uint64 {
	if l.CPU > 0 && v > l.CPU {
		return l.CPU
	}
	return v
}
