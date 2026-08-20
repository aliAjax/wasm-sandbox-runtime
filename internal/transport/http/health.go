package httptransport

import (
	"net/http"
	"sync/atomic"
)

type Health struct {
	ready    atomic.Bool
	draining atomic.Bool
}

func (h *Health) SetReady(v bool) { h.ready.Store(v) }
func (h *Health) Drain()          { h.draining.Store(true); h.ready.Store(false) }
func (h *Health) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/readyz" && !h.ready.Load() {
		w.WriteHeader(503)
		return
	}
	if h.draining.Load() {
		w.WriteHeader(503)
		return
	}
	jsonWrite(w, 200, map[string]string{"status": "ok"})
}
