package application

import (
	"sync"
	"time"
)

type Usage struct {
	TenantID        string    `json:"tenant_id"`
	Executions      uint64    `json:"executions"`
	CPUInstructions uint64    `json:"cpu_instructions"`
	MemoryPeak      uint64    `json:"memory_peak"`
	OutputBytes     uint64    `json:"output_bytes"`
	UpdatedAt       time.Time `json:"updated_at"`
}
type UsageBook struct {
	mu    sync.Mutex
	items map[string]Usage
}

func NewUsageBook() *UsageBook { return &UsageBook{items: map[string]Usage{}} }
func (u *UsageBook) Record(tenant string, instructions, memory, output uint64, at time.Time) {
	u.mu.Lock()
	defer u.mu.Unlock()
	v := u.items[tenant]
	v.TenantID = tenant
	v.Executions++
	v.CPUInstructions += instructions
	if memory > v.MemoryPeak {
		v.MemoryPeak = memory
	}
	v.OutputBytes += output
	v.UpdatedAt = at
	u.items[tenant] = v
}
func (u *UsageBook) Get(tenant string) Usage {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.items[tenant]
}
func (u *UsageBook) Snapshot() []Usage {
	u.mu.Lock()
	defer u.mu.Unlock()
	out := make([]Usage, 0, len(u.items))
	for _, v := range u.items {
		out = append(out, v)
	}
	return out
}
