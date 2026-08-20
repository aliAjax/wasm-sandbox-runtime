package worker

import (
	"sync"
	"time"
)

type WorkerHealth struct {
	ID            string
	StartedAt     time.Time
	LastHeartbeat time.Time
	Active        int
	Capacity      int
	Version       string
}
type Registry struct {
	mu      sync.RWMutex
	workers map[string]WorkerHealth
}

func NewRegistry() *Registry { return &Registry{workers: map[string]WorkerHealth{}} }
func (r *Registry) Register(id, version string, capacity int, now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workers[id] = WorkerHealth{ID: id, Version: version, Capacity: capacity, StartedAt: now, LastHeartbeat: now}
}
func (r *Registry) Heartbeat(id string, active int, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.workers[id]
	if !ok {
		return false
	}
	w.Active = active
	w.LastHeartbeat = now
	r.workers[id] = w
	return true
}
func (r *Registry) Healthy(now time.Time, stale time.Duration) []WorkerHealth {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []WorkerHealth{}
	for _, w := range r.workers {
		if now.Sub(w.LastHeartbeat) <= stale {
			out = append(out, w)
		}
	}
	return out
}
