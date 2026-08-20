package quota

import "errors"

func (q Quota) Available() (int, float64, uint64, uint64) {
	return q.MaxConcurrent - q.UsedConcurrent, q.CPUSeconds - q.UsedCPU, q.MemoryBytes - q.UsedMemory, q.StorageBytes - q.UsedStorage
}
func (q Quota) CanReserve(cpu float64, memory, storage uint64) error {
	c, cu, cm, cs := q.Available()
	if c < 1 {
		return errors.New("concurrency unavailable")
	}
	if cu < cpu {
		return errors.New("cpu unavailable")
	}
	if cm < memory {
		return errors.New("memory unavailable")
	}
	if cs < storage {
		return errors.New("storage unavailable")
	}
	return nil
}
func (q Quota) Utilization() (float64, float64, float64) {
	cp := 0.0
	if q.CPUSeconds > 0 {
		cp = q.UsedCPU / q.CPUSeconds
	}
	mem := 0.0
	if q.MemoryBytes > 0 {
		mem = float64(q.UsedMemory) / float64(q.MemoryBytes)
	}
	store := 0.0
	if q.StorageBytes > 0 {
		store = float64(q.UsedStorage) / float64(q.StorageBytes)
	}
	return cp, mem, store
}
