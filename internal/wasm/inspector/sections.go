package inspector

import "errors"

type SectionID byte

const (
	Custom          SectionID = 0
	TypeSection     SectionID = 1
	ImportSection   SectionID = 2
	FunctionSection SectionID = 3
	MemorySection   SectionID = 5
	ExportSection   SectionID = 7
	CodeSection     SectionID = 10
)

func CheckMemory(info ModuleInfo, maxPages uint32) error {
	if info.MemoryPages > maxPages {
		return errors.New("linear memory limit exceeded")
	}
	return nil
}
func CheckInstructionBudget(info ModuleInfo, max uint64) error {
	if info.EstimatedInstructions > max {
		return errors.New("instruction budget exceeded")
	}
	return nil
}
