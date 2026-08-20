package module

func (r ResourceLimits) CPUSeconds() float64 { return float64(r.CPUInstructions) / 1e9 }
func (r ResourceLimits) MemoryMiB() uint64   { return r.MemoryBytes / (1024 * 1024) }
func (r ResourceLimits) OutputMiB() uint64   { return r.OutputBytes / (1024 * 1024) }
func (r ResourceLimits) Tightened(other ResourceLimits) ResourceLimits {
	out := r
	if other.CPUInstructions > 0 && other.CPUInstructions < out.CPUInstructions {
		out.CPUInstructions = other.CPUInstructions
	}
	if other.MemoryBytes > 0 && other.MemoryBytes < out.MemoryBytes {
		out.MemoryBytes = other.MemoryBytes
	}
	if other.OutputBytes > 0 && other.OutputBytes < out.OutputBytes {
		out.OutputBytes = other.OutputBytes
	}
	return out
}
func (r ResourceLimits) IsTightenedBy(other ResourceLimits) bool { return other.Tightened(r) != r }
