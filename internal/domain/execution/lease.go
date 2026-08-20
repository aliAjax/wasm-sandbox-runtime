package execution

import (
	"errors"
	"time"
)

type Lease struct {
	ID          string
	ExecutionID string
	WorkerID    string
	ExpiresAt   time.Time
	Version     int64
}

func (l Lease) Valid(now time.Time, version int64) error {
	if l.ID == "" || l.ExecutionID == "" || l.WorkerID == "" {
		return errors.New("lease identity missing")
	}
	if !now.Before(l.ExpiresAt) {
		return errors.New("lease expired")
	}
	if version != l.Version {
		return errors.New("stale lease version")
	}
	return nil
}
func (l *Lease) Renew(now time.Time, d time.Duration) { l.ExpiresAt = now.Add(d); l.Version++ }
func (l Lease) Remaining(now time.Time) time.Duration {
	if now.After(l.ExpiresAt) {
		return 0
	}
	return l.ExpiresAt.Sub(now)
}
