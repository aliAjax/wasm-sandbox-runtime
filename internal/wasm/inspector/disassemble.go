package inspector

import "fmt"

func Disassemble(code []byte) []string {
	out := make([]string, 0, len(code))
	for pc, op := range code {
		out = append(out, fmt.Sprintf("%04d %s", pc, OpcodeName(op)))
	}
	return out
}
func InstructionBudget(code []byte, perByte uint64) uint64 { return uint64(len(code))*perByte + 1 }
