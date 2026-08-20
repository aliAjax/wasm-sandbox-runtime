package policy

type Change struct {
	Field  string
	Before any
	After  any
	Risk   string
}

func Diff(a, b Policy) []Change {
	out := []Change{}
	if a.Limits != b.Limits {
		out = append(out, Change{Field: "limits", Before: a.Limits, After: b.Limits, Risk: "high"})
	}
	if a.Deterministic != b.Deterministic {
		out = append(out, Change{Field: "deterministic", Before: a.Deterministic, After: b.Deterministic, Risk: "medium"})
	}
	if len(a.Network.Domains) != len(b.Network.Domains) {
		out = append(out, Change{Field: "network.domains", Before: a.Network.Domains, After: b.Network.Domains, Risk: "high"})
	}
	return out
}
func Breaking(changes []Change) bool {
	for _, c := range changes {
		if c.Risk == "high" {
			return true
		}
	}
	return false
}
