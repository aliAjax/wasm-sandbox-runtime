package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/application"
)

type job struct {
	id       string
	priority int
	created  time.Time
}
type Scheduler struct {
	service  *application.Service
	workers  int
	queue    chan job
	stop     chan struct{}
	done     chan struct{}
	mu       sync.Mutex
	stopOnce sync.Once
	seen     map[string]struct{}
	log      *slog.Logger
}

func NewScheduler(s *application.Service, n int) *Scheduler {
	if n < 1 {
		n = 1
	}
	return &Scheduler{service: s, workers: n, queue: make(chan job, 256), stop: make(chan struct{}), done: make(chan struct{}), seen: map[string]struct{}{}, log: slog.Default()}
}
func (s *Scheduler) Start() {
	go s.dispatch()
	for i := 0; i < s.workers; i++ {
		go s.run(i)
	}
}
func (s *Scheduler) dispatch() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	defer close(s.done)
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.enqueue()
		}
	}
}
func (s *Scheduler) enqueue() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	list, err := s.service.ListExecutions(ctx, "", 0, 256)
	if err != nil {
		return
	}
	for _, e := range list {
		if e.State != "queued" {
			continue
		}
		s.mu.Lock()
		_, exists := s.seen[e.ID]
		if !exists {
			s.seen[e.ID] = struct{}{}
		}
		s.mu.Unlock()
		if !exists {
			select {
			case s.queue <- job{id: e.ID, priority: e.Priority, created: e.CreatedAt}:
			default:
				return
			}
		}
	}
}
func (s *Scheduler) run(index int) {
	workerID := fmt.Sprintf("worker-%d", index)
	for {
		select {
		case <-s.stop:
			return
		case j := <-s.queue:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			if err := s.service.Admit(ctx, j.id); err == nil {
				_ = s.service.Run(ctx, j.id, workerID)
			}
			cancel()
			s.mu.Lock()
			delete(s.seen, j.id)
			s.mu.Unlock()
		}
	}
}
func (s *Scheduler) Stop() {
	close(s.stop)
	select {
	case <-s.done:
	case <-time.After(time.Second):
	}
}
