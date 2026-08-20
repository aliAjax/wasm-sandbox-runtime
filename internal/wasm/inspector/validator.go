package inspector

import (
	"errors"
	"fmt"
)

type Limits struct {
	MaxMemoryPages  uint32
	MaxInstructions uint64
	MaxImports      int
	MaxExports      int
}

func Validate(info ModuleInfo, l Limits) error {
	if len(info.Imports) > l.MaxImports {
		return errors.New("too many imports")
	}
	if len(info.Exports) > l.MaxExports {
		return errors.New("too many exports")
	}
	if err := CheckMemory(info, l.MaxMemoryPages); err != nil {
		return err
	}
	return CheckInstructionBudget(info, l.MaxInstructions)
}
func RequiredExports(info ModuleInfo) map[string]bool {
	out := map[string]bool{}
	for _, e := range info.Exports {
		if e.Kind == 0 {
			out[e.Name] = true
		}
	}
	return out
}
func ImportNames(info ModuleInfo) []string {
	out := make([]string, 0, len(info.Imports))
	for _, v := range info.Imports {
		out = append(out, fmt.Sprintf("%s.%s", v.Module, v.Name))
	}
	return out
}
