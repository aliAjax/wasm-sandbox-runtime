package observability

import (
	"sync/atomic"
	"time"
)

type State struct {
	ready   atomic.Bool
	started time.Time
}

func NewHealth() *State                { h := &State{started: time.Now().UTC()}; h.ready.Store(true); return h }
func (h *State) Ready() bool           { return h.ready.Load() }
func (h *State) Stop()                 { h.ready.Store(false) }
func (h *State) Uptime() time.Duration { return time.Since(h.started) }
