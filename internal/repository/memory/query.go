package memory

import (
	"sort"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
)

func sortExecutions(items []execution.Execution) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		return items[i].Priority > items[j].Priority
	})
}
func stateCount(items []execution.Execution) map[execution.State]int {
	out := map[execution.State]int{}
	for _, e := range items {
		out[e.State]++
	}
	return out
}
func since(items []execution.Execution, t time.Time) []execution.Execution {
	out := []execution.Execution{}
	for _, e := range items {
		if e.CreatedAt.After(t) {
			out = append(out, e)
		}
	}
	return out
}
