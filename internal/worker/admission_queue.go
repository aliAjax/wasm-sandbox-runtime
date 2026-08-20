package worker

import (
	"sync"
)

type AdmissionQueue struct {
	mu       sync.Mutex
	items    []job
	capacity int
	closed   bool
}

func (q *AdmissionQueue) Close() { q.mu.Lock(); q.closed = true; q.mu.Unlock() }

func NewAdmissionQueue(capacity int) *AdmissionQueue {
	if capacity < 1 {
		capacity = 128
	}
	return &AdmissionQueue{capacity: capacity}
}
func (q *AdmissionQueue) Push(j job) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	if len(q.items) >= q.capacity {
		return false
	}
	q.items = append(q.items, j)
	return true
}
func (q *AdmissionQueue) Pop() (job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return job{}, false
	}
	j := q.items[0]
	q.items = q.items[1:]
	return j, true
}
func (q *AdmissionQueue) Size() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
