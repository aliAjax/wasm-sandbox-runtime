package execution

import "time"

type Filter struct {
	TenantID string
	States   []State
	Since    *time.Time
	Until    *time.Time
}

func (f Filter) Match(e Execution) bool {
	if f.TenantID != "" && f.TenantID != e.TenantID {
		return false
	}
	if len(f.States) > 0 {
		ok := false
		for _, s := range f.States {
			if s == e.State {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	if f.Since != nil && e.CreatedAt.Before(*f.Since) {
		return false
	}
	if f.Until != nil && !e.CreatedAt.Before(*f.Until) {
		return false
	}
	return true
}
func FilterAll(items []Execution, f Filter) []Execution {
	out := []Execution{}
	for _, e := range items {
		if f.Match(e) {
			out = append(out, e)
		}
	}
	return out
}
