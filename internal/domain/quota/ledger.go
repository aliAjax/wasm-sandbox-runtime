package quota

import (
	"errors"
	"sync"
	"time"
)

type Entry struct {
	ID       string
	TenantID string
	Kind     string
	CPU      float64
	Memory   uint64
	Storage  uint64
	At       time.Time
}
type Ledger struct {
	mu      sync.Mutex
	entries []Entry
}

func (l *Ledger) Append(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e.ID == "" || e.TenantID == "" {
		return errors.New("ledger entry identity missing")
	}
	for _, old := range l.entries {
		if old.ID == e.ID {
			return errors.New("duplicate ledger entry")
		}
	}
	l.entries = append(l.entries, e)
	return nil
}
func (l *Ledger) List(tenant string, from, to time.Time) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := []Entry{}
	for _, e := range l.entries {
		if e.TenantID == tenant && (from.IsZero() || !e.At.Before(from)) && (to.IsZero() || e.At.Before(to)) {
			out = append(out, e)
		}
	}
	return out
}
func (l *Ledger) Totals(tenant string) (cpu float64, memory, storage uint64) {
	for _, e := range l.List(tenant, time.Time{}, time.Time{}) {
		cpu += e.CPU
		memory += e.Memory
		storage += e.Storage
	}
	return
}
