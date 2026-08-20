package artifact

import "time"

type Retention struct {
	ExpiresAt time.Time
	LegalHold bool
	Deleted   bool
}

func (r Retention) Expired(now time.Time) bool {
	return !r.LegalHold && !r.Deleted && !r.ExpiresAt.IsZero() && now.After(r.ExpiresAt)
}
func (r *Retention) Extend(until time.Time) {
	if until.After(r.ExpiresAt) {
		r.ExpiresAt = until
	}
}
func (r *Retention) MarkDeleted() {
	if !r.LegalHold {
		r.Deleted = true
	}
}
