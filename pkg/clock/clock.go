package clock

import "time"

type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}
type System struct{}

func (System) Now() time.Time                         { return time.Now().UTC() }
func (System) After(d time.Duration) <-chan time.Time { return time.After(d) }

type Fixed struct{ Value time.Time }

func (f Fixed) Now() time.Time                         { return f.Value }
func (f Fixed) After(d time.Duration) <-chan time.Time { return time.After(d) }
