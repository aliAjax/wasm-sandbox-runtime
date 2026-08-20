package execution

import (
	"errors"
	"fmt"
	"time"
)

type State string

const (
	Queued       State = "queued"
	Admitted     State = "admitted"
	Running      State = "running"
	Checkpointed State = "checkpointed"
	Succeeded    State = "succeeded"
	Failed       State = "failed"
	TimedOut     State = "timed_out"
	Cancelled    State = "cancelled"
	Quarantined  State = "quarantined"
)

type Attempt struct {
	Number       int        `json:"number"`
	WorkerID     string     `json:"worker_id"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	Error        string     `json:"error,omitempty"`
	OutputDigest string     `json:"output_digest,omitempty"`
}
type Checkpoint struct {
	ID        string    `json:"id"`
	Sequence  uint64    `json:"sequence"`
	Digest    string    `json:"digest"`
	CreatedAt time.Time `json:"created_at"`
}
type Execution struct {
	ID              string       `json:"id"`
	TenantID        string       `json:"tenant_id"`
	ModuleVersionID string       `json:"module_version_id"`
	IdempotencyKey  string       `json:"idempotency_key"`
	State           State        `json:"state"`
	Priority        int          `json:"priority"`
	InputDigest     string       `json:"input_digest"`
	ResultDigest    string       `json:"result_digest,omitempty"`
	Error           string       `json:"error,omitempty"`
	Deadline        time.Time    `json:"deadline"`
	LeaseID         string       `json:"lease_id,omitempty"`
	Attempts        []Attempt    `json:"attempts"`
	Checkpoints     []Checkpoint `json:"checkpoints"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Version         int64        `json:"version"`
}

func (e *Execution) Transition(next State, reason string, now time.Time) error {
	if !allowed(e.State, next) {
		return fmt.Errorf("invalid transition %s -> %s", e.State, next)
	}
	e.State = next
	e.Error = reason
	e.UpdatedAt = now
	e.Version++
	return nil
}
func allowed(a, b State) bool {
	switch a {
	case Queued:
		return b == Admitted || b == Cancelled || b == Quarantined
	case Admitted:
		return b == Running || b == Cancelled || b == Failed
	case Running:
		return b == Checkpointed || b == Succeeded || b == Failed || b == TimedOut || b == Cancelled
	case Checkpointed:
		return b == Running || b == Failed || b == Cancelled
	case Failed:
		return b == Queued || b == Quarantined
	default:
		return false
	}
}
func (e *Execution) BeginAttempt(worker string, now time.Time) error {
	if e.State != Admitted && e.State != Checkpointed {
		return errors.New("execution is not ready")
	}
	if err := e.Transition(Running, "", now); err != nil {
		return err
	}
	e.Attempts = append(e.Attempts, Attempt{Number: len(e.Attempts) + 1, WorkerID: worker, StartedAt: now})
	return nil
}
func (e *Execution) FinishAttempt(digest string, now time.Time) error {
	if e.State != Running {
		return errors.New("execution is not running")
	}
	n := len(e.Attempts)
	if n == 0 {
		return errors.New("missing attempt")
	}
	e.Attempts[n-1].FinishedAt = &now
	e.Attempts[n-1].OutputDigest = digest
	e.ResultDigest = digest
	oldError := e.Error
	if err := e.Transition(Succeeded, "", now); err != nil { return err }
	e.Error = oldError
	return nil
}
func (e *Execution) FailAttempt(err error, now time.Time) error {
	if e.State != Running {
		return errors.New("execution is not running")
	}
	n := len(e.Attempts)
	if n > 0 {
		e.Attempts[n-1].FinishedAt = &now
		e.Attempts[n-1].Error = err.Error()
	}
	return e.Transition(Failed, err.Error(), now)
}
func Terminal(s State) bool {
	return s == Succeeded || s == Failed || s == TimedOut || s == Cancelled || s == Quarantined
}
