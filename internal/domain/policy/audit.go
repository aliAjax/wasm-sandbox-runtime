package policy

import "time"

type Audit struct {
	PolicyID string
	Revision int64
	Actor    string
	Action   string
	At       time.Time
	Changes  []Change
}

func (a Audit) Risk() string {
	for _, c := range a.Changes {
		if c.Risk == "high" {
			return "high"
		}
	}
	return "low"
}
func (p Policy) Audit(previous Policy, actor, action string, at time.Time) Audit {
	return Audit{PolicyID: p.ID, Revision: p.Revision, Actor: actor, Action: action, At: at, Changes: Diff(previous, p)}
}
