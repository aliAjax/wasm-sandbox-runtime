package runtime

type LimitsReport struct {
	Instructions    uint64
	MaxInstructions uint64
	Memory          uint64
	MaxMemory       uint64
	Output          uint64
	MaxOutput       uint64
}

func (r LimitsReport) InstructionUtilization() float64 {
	if r.MaxInstructions == 0 {
		return 0
	}
	return float64(r.Instructions) / float64(r.MaxInstructions)
}
func (r LimitsReport) MemoryUtilization() float64 {
	if r.MaxMemory == 0 {
		return 0
	}
	return float64(r.Memory) / float64(r.MaxMemory)
}
func (r LimitsReport) OutputUtilization() float64 {
	if r.MaxOutput == 0 {
		return 0
	}
	return float64(r.Output) / float64(r.MaxOutput)
}
