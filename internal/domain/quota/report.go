package quota

import "time"

type Report struct {
	TenantID       string
	At             time.Time
	ConcurrentUsed int
	ConcurrentMax  int
	CPUUsed        float64
	CPUMax         float64
	MemoryUsed     uint64
	MemoryMax      uint64
	StorageUsed    uint64
	StorageMax     uint64
}

func BuildReport(q Quota, at time.Time) Report {
	return Report{TenantID: q.TenantID, At: at, ConcurrentUsed: q.UsedConcurrent, ConcurrentMax: q.MaxConcurrent, CPUUsed: q.UsedCPU, CPUMax: q.CPUSeconds, MemoryUsed: q.UsedMemory, MemoryMax: q.MemoryBytes, StorageUsed: q.UsedStorage, StorageMax: q.StorageBytes}
}
func (r Report) Saturated() bool {
	return r.ConcurrentUsed >= r.ConcurrentMax || r.CPUUsed >= r.CPUMax || r.MemoryUsed >= r.MemoryMax || r.StorageUsed >= r.StorageMax
}
