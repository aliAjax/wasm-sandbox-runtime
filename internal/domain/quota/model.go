package quota

import (
	"errors"
	"sync"
	"time"
)

type Quota struct {
	TenantID       string    `json:"tenant_id"`
	MaxConcurrent  int       `json:"max_concurrent"`
	CPUSeconds     float64   `json:"cpu_seconds"`
	MemoryBytes    uint64    `json:"memory_bytes"`
	StorageBytes   uint64    `json:"storage_bytes"`
	UsedConcurrent int       `json:"used_concurrent"`
	UsedCPU        float64   `json:"used_cpu"`
	UsedMemory     uint64    `json:"used_memory"`
	UsedStorage    uint64    `json:"used_storage"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type Reservation struct {
	TenantID  string
	CPU       float64
	Memory    uint64
	Storage   uint64
	CreatedAt time.Time
	released  bool
	mu        sync.Mutex
}

func (q *Quota) Reserve(cpu float64, memory, storage uint64, now time.Time) (*Reservation, error) {
	if q.UsedConcurrent >= q.MaxConcurrent {
		return nil, errors.New("concurrent quota exhausted")
	}
	if q.UsedCPU+cpu > q.CPUSeconds {
		return nil, errors.New("cpu quota exhausted")
	}
	if q.UsedMemory+memory > q.MemoryBytes {
		return nil, errors.New("memory quota exhausted")
	}
	if q.UsedStorage+storage > q.StorageBytes {
		return nil, errors.New("storage quota exhausted")
	}
	q.UsedConcurrent++
	q.UsedCPU += cpu
	q.UsedMemory += memory
	q.UsedStorage += storage
	q.UpdatedAt = now
	return &Reservation{TenantID: q.TenantID, CPU: cpu, Memory: memory, Storage: storage, CreatedAt: now}, nil
}
func (q *Quota) Release(r *Reservation, now time.Time) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.released {
		return
	}
	r.released = true
	if q.UsedConcurrent > 0 {
		q.UsedConcurrent--
	}
	if q.UsedCPU >= r.CPU {
		q.UsedCPU -= r.CPU
	} else {
		q.UsedCPU = 0
	}
	if q.UsedMemory >= r.Memory {
		q.UsedMemory -= r.Memory
	} else {
		q.UsedMemory = 0
	}
	if q.UsedStorage >= r.Storage {
		q.UsedStorage -= r.Storage
	} else {
		q.UsedStorage = 0
	}
	q.UpdatedAt = now
}
func (q Quota) Validate() error {
	if q.TenantID == "" || q.MaxConcurrent < 1 || q.CPUSeconds <= 0 || q.MemoryBytes < 64*1024 || q.StorageBytes < 1 {
		return errors.New("invalid quota")
	}
	return nil
}
