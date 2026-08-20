package worker

import (
	"context"
	"errors"
	"testing"
)

func TestSchedulerStopDrainsWorkers(t *testing.T) {
	s := NewScheduler(nil, 1)
	s.Start()
	s.Stop()
	panicked := false
	func() { defer func() { panicked = recover() != nil }(); s.Stop() }()
	if panicked {
		t.Fatal("Stop is not idempotent")
	}
	q := NewAdmissionQueue(1)
	q.Close()
	start := make(chan struct{})
	results := make(chan bool, 2)
	go func() { <-start; results <- q.Push(job{id: "late-1"}) }()
	go func() { <-start; results <- q.Push(job{id: "late-2"}) }()
	close(start)
	if <-results || <-results {
		t.Fatal("closed queue accepted work")
	}
}

func TestRecoveryStopsOnContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if (Recovery{MaxAttempts: 3}).ShouldRetryContext(ctx, 1, errors.New("retry")) {
		t.Fatal("cancelled recovery continued")
	}
}
func TestPriorityQueueEmptyIsSafe(t *testing.T) {
	q := PriorityQueue{}
	if _, ok := q.TakeSafe(); ok {
		t.Fatal("empty priority queue returned a job")
	}
}
