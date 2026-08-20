package httptransport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SSE struct {
	w      http.ResponseWriter
	f      http.Flusher
	closed bool
}

func NewSSE(w http.ResponseWriter) (*SSE, bool) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, false
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	return &SSE{w: w, f: f}, true
}
func (s *SSE) Send(event string, v any) {
	if s.closed {
		return
	}
	b, _ := json.Marshal(v)
	fmt.Fprintf(s.w, "id: %d\nevent: %s\ndata: %s\n\n", time.Now().UnixNano(), event, b)
	s.f.Flush()
}
func (s *SSE) Close() { s.closed = true }
