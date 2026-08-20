package execution

import (
	"errors"
	"time"
)

type TransitionEvent struct {
	ExecutionID string
	From        State
	To          State
	Reason      string
	At          time.Time
	Version     int64
}

func (e Execution) CanCancel() bool {
	return e.State == Queued || e.State == Admitted || e.State == Running || e.State == Checkpointed
}
func (e Execution) CanCheckpoint() bool { return e.State == Running }
func (e Execution) CurrentAttempt() (*Attempt, error) {
	if len(e.Attempts) == 0 {
		return nil, errors.New("no attempts")
	}
	return &e.Attempts[len(e.Attempts)-1], nil
}
func (e *Execution) Quarantine(reason string, at time.Time) error {
	if Terminal(e.State) {
		return errors.New("terminal execution")
	}
	e.State = Quarantined
	e.Error = reason
	e.UpdatedAt = at
	e.Version++
	return nil
}
func (e *Execution) Retry(at time.Time) error {
	if e.State != Failed {
		return errors.New("only failed execution can retry")
	}
	e.State = Queued
	e.Error = ""
	e.UpdatedAt = at
	e.Version++
	return nil
}
func (e Execution) Duration() time.Duration {
	if len(e.Attempts) == 0 {
		return 0
	}
	a := e.Attempts[0]
	end := e.UpdatedAt
	if a.FinishedAt != nil {
		end = *a.FinishedAt
	}
	return end.Sub(a.StartedAt)
}
func (e Execution) AttemptCount() int { return len(e.Attempts) }
