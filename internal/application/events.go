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
		subs := h.subscribers[id]
		if subs != nil {
			delete(subs, ch)
			if len(subs) == 0 {
				delete(h.subscribers, id)
			}
		}
		h.mu.Unlock()
		close(ch)
	}
}
func (h *EventHub) Publish(e Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers[e.ExecutionID] {
		select {
		case ch <- e:
		default:
		}
	}
}
