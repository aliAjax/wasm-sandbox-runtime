package worker

import "time"

type Backoff struct {
	Base     time.Duration
	Max      time.Duration
	Attempts int
}

func (b Backoff) Duration(n int) time.Duration {
	if n < 0 {
		n = 0
	}
	d := b.Base
	for i := 0; i < n; i++ {
		d *= 2
		if d >= b.Max {
			return b.Max
		}
	}
	if d > b.Max {
		return b.Max
	}
	return d
}

type Lease struct {
	ID        string
	Worker    string
	ExpiresAt time.Time
	Version   int64
}

func (l Lease) Valid(now time.Time) bool { return l.ID != "" && now.Before(l.ExpiresAt) }
