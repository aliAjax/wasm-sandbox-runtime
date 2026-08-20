package application

import (
	"sync"
	"time"
)

type Event struct {
	ExecutionID  string            `json:"execution_id"`
	State        string            `json:"state"`
	ResultDigest string            `json:"result_digest,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
	At           time.Time         `json:"at"`
}
type EventHub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan Event]struct{}
}

func NewEventHub() *EventHub { return &EventHub{subscribers: map[string]map[chan Event]struct{}{}} }
func (h *EventHub) Subscribe(id string) (<-chan Event, func()) {
	ch := make(chan Event, 16)
	h.mu.Lock()
	if h.subscribers[id] == nil {
		h.subscribers[id] = map[chan Event]struct{}{}
	}
	h.subscribers[id][ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		h.mu.Unlock()
		close(ch)
	}
}
func (h *EventHub) Publish(e Event) {
	h.mu.RLock()
	channels := make([]chan Event, 0, len(h.subscribers[e.ExecutionID]))
	for ch := range h.subscribers[e.ExecutionID] {
		channels = append(channels, ch)
	}
	h.mu.RUnlock()
	for _, ch := range channels {
		select {
		case ch <- e:
		default:
		}
	}
}
