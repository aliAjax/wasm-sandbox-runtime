package inspector

import "strings"

func Risk(info ModuleInfo) string {
	score := 0
	for _, im := range info.Imports {
		if strings.Contains(im.Module, "wasi") {
			score++
		}
		if strings.Contains(im.Name, "path") || strings.Contains(im.Name, "socket") {
			score += 3
		}
	}
	if info.MemoryPages > 16384 {
		score += 2
	}
	if score >= 6 {
		return "high"
	}
	if score >= 2 {
		return "medium"
	}
	return "low"
}
func Safe(info ModuleInfo) bool { return Risk(info) != "high" }
