package module

import "strings"

func NormalizeImports(in []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, v := range in {
		v = strings.TrimSpace(strings.ToLower(v))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
func DangerousImports(in []string) []string {
	out := []string{}
	for _, v := range NormalizeImports(in) {
		if strings.Contains(v, "wasi_snapshot_preview1.path_open") || strings.Contains(v, "network") || strings.Contains(v, "proc_exit") {
			out = append(out, v)
		}
	}
	return out
}
