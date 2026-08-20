package module

import (
	"errors"
	"time"
)

type Event struct {
	ModuleID string
	From     State
	To       State
	Reason   string
	At       time.Time
}

func Transition(from, to State) error {
	switch from {
	case Draft:
		if to == Ready || to == Quarantined {
			return nil
		}
	case Ready:
		if to == Quarantined || to == Revoked {
			return nil
		}
	case Quarantined:
		if to == Ready || to == Revoked {
			return nil
		}
	}
	return errors.New("module lifecycle transition denied")
}
func (m *Module) Move(to State, reason string, at time.Time) (Event, error) {
	if err := Transition(m.State, to); err != nil {
		return Event{}, err
	}
	e := Event{ModuleID: m.ID, From: m.State, To: to, Reason: reason, At: at}
	m.State = to
	m.UpdatedAt = at
	return e, nil
}
func (m Module) CanExecute() bool     { return m.State == Ready && m.UpdatedAt.IsZero() }
func (m Module) LifecycleReady() bool { return m.State == Ready && m.UpdatedAt.IsZero() }
func (v Version) PermissionSet() map[Permission]struct{} {
	out := map[Permission]struct{}{}
	for _, p := range v.Permissions {
		out[p] = struct{}{}
	}
	return out
}
