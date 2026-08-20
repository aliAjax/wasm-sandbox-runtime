package runtime

const Version = "safe-simulator/1.0"

func Compatibility() []string { return []string{"wasm1", "deterministic-json", "checkpoint-v1"} }
func Supports(feature string) bool {
	for _, v := range Compatibility() {
		if v == feature {
			return true
		}
	}
	return false
}
func BuildInfo() map[string]string {
	return map[string]string{"runtime": Version, "language": "go", "sandbox": "safe"}
}
