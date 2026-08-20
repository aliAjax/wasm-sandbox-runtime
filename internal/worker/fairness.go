package worker

import (
	"sort"
	"sync"
	"time"
)

type TenantQueue struct {
	Tenant     string
	Weight     int
	LastServed time.Time
	Jobs       []job
}
type FairQueue struct {
	mu      sync.Mutex
	tenants map[string]*TenantQueue
}

func NewFairQueue() *FairQueue { return &FairQueue{tenants: map[string]*TenantQueue{}} }
func (q *FairQueue) Add(tenant string, j job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	t := q.tenants[tenant]
	if t == nil {
		t = &TenantQueue{Tenant: tenant, Weight: 1}
		q.tenants[tenant] = t
	}
	t.Jobs = append(t.Jobs, j)
}
func (q *FairQueue) Take() (job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	list := make([]*TenantQueue, 0, len(q.tenants))
	for _, t := range q.tenants {
		if len(t.Jobs) > 0 {
			list = append(list, t)
		}
	}
	if len(list) == 0 {
		return job{}, false
	}
	sort.Slice(list, func(i, j int) bool { return list[i].LastServed.Before(list[j].LastServed) })
	t := list[0]
	v := t.Jobs[0]
	t.Jobs = t.Jobs[1:]
	t.LastServed = time.Now()
	return v, true
}
func (q *FairQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := 0
	for _, t := range q.tenants {
		n += len(t.Jobs)
	}
	return n
}
