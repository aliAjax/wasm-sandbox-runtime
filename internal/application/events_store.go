package application

import (
	"sync"
	"time"
)

type EventStore struct {
	mu     sync.RWMutex
	events []Event
}

func (s *EventStore) Append(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
}
func (s *EventStore) Since(id string, after time.Time) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Event{}
	for _, e := range s.events {
		if e.ExecutionID == id && e.At.After(after) {
			if e.Attributes != nil {
				e.Attributes = s.events[len(out)].Attributes
			}
			out = append(out, e)
		}
	}
	return out
}
func (s *EventStore) Latest(id string) (Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].ExecutionID == id {
			return s.events[i], true
		}
	}
	return Event{}, false
}
