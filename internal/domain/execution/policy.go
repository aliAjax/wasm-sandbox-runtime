package execution

import "errors"

type ExecutionPolicy struct {
	MaxAttempts        int
	AllowCheckpoint    bool
	AllowRetry         bool
	CancelGraceSeconds int
}

func (p ExecutionPolicy) Validate() error {
	if p.MaxAttempts < 1 || p.MaxAttempts > 10 {
		return errors.New("max attempts invalid")
	}
	if p.CancelGraceSeconds < 0 || p.CancelGraceSeconds > 300 {
		return errors.New("cancel grace invalid")
	}
	return nil
}
func (p ExecutionPolicy) CanRetry(attempt int) bool { return p.AllowRetry && attempt < p.MaxAttempts }
func (p ExecutionPolicy) CanCheckpoint() bool       { return p.AllowCheckpoint }
