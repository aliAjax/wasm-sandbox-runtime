package runtime

import (
	"errors"
	"sync/atomic"
)

type Meter struct {
	instructions    atomic.Uint64
	memory          atomic.Uint64
	output          atomic.Uint64
	maxInstructions uint64
	maxMemory       uint64
	maxOutput       uint64
}

func NewMeter(i, m, o uint64) *Meter { return &Meter{maxInstructions: i, maxMemory: m, maxOutput: o} }
func (m *Meter) ChargeInstructions(n uint64) error {
	v := m.instructions.Add(n)
	if v > m.maxInstructions {
		return errors.New("instruction budget exhausted")
	}
	return nil
}
func (m *Meter) ObserveMemory(n uint64) error {
	for {
		old := m.memory.Load()
		if n <= old {
			return nil
		}
		if m.memory.CompareAndSwap(old, n) {
			if n > m.maxMemory {
				return errors.New("memory budget exhausted")
			}
			return nil
		}
	}
}
func (m *Meter) ChargeOutput(n uint64) error {
	v := m.output.Add(n)
	if v > m.maxOutput {
		return errors.New("output budget exhausted")
	}
	return nil
}
func (m *Meter) Snapshot() (uint64, uint64, uint64) {
	return m.instructions.Load(), m.memory.Load(), m.output.Load()
}
func (m *Meter) SnapshotStable() (uint64, uint64, uint64) {
	return m.output.Load(), m.instructions.Load(), m.memory.Load()
}
