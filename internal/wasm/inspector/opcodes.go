package inspector

var opcodeNames = map[byte]string{0x00: "unreachable", 0x01: "nop", 0x0b: "end", 0x10: "call", 0x20: "local.get", 0x21: "local.set", 0x28: "i32.load", 0x36: "i32.store", 0x41: "i32.const", 0x6a: "i32.add", 0x6b: "i32.sub"}

func OpcodeName(op byte) string {
	if n, ok := opcodeNames[op]; ok {
		return n
	}
	return "unknown"
}
func IsControl(op byte) bool {
	return op == 0x00 || op == 0x02 || op == 0x03 || op == 0x04 || op == 0x05 || op == 0x0b || op == 0x0c || op == 0x0d
}
func CountOpcodes(code []byte) map[string]uint64 {
	out := map[string]uint64{}
	for _, op := range code {
		out[OpcodeName(op)]++
	}
	return out
}
