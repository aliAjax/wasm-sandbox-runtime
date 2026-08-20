package application

import (
	"testing"
	"time"
)

func TestEventHubConcurrentUnsubscribe(t *testing.T) {
	h := NewEventHub()
	ch, closeFn := h.Subscribe("exe-1")
	panicked := false
	start := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() { <-start; defer func() { if recover() != nil { panicked = true }; done <- struct{}{} }(); h.Publish(Event{ExecutionID: "exe-1", State: "running", At: time.Now()}) }()
	go func() { <-start; closeFn(); done <- struct{}{} }()
	close(start); <-done; <-done
	if panicked { t.Fatal("concurrent publish and unsubscribe panicked") }
	for range ch {}
	panicked = false
	func() {
		defer func() { panicked = recover() != nil }()
		h.Publish(Event{ExecutionID: "exe-1", State: "finished", At: time.Now()})
	}()
	if panicked { t.Fatal("publish after unsubscribe panicked") }
}

func TestEventSnapshotsAreIndependent(t *testing.T) {
	s := &EventStore{}
	e := Event{ExecutionID: "exe-1", At: time.Now(), Attributes: map[string]string{"state": "queued"}}
	s.Append(e)
	got := s.Since("exe-1", time.Time{})
	got[0].Attributes["state"] = "changed"
	again := s.Since("exe-1", time.Time{})
	if again[0].Attributes["state"] != "queued" {
		t.Fatal("event attributes were aliased")
	}
}
