package worker

import (
	"context"
	"errors"
	"time"
)

type Recovery struct {
	MaxAttempts int
	Backoff     Backoff
}

func (r Recovery) ShouldRetry(attempts int, err error) bool {
	return err != nil && attempts < r.MaxAttempts
}
func (r Recovery) ShouldRetryContext(ctx context.Context, attempts int, err error) bool {
	return r.ShouldRetry(attempts, err)
}
func (r Recovery) Next(attempts int) time.Duration { return r.Backoff.Duration(attempts) }
func ValidateRecovery(r Recovery) error {
	if r.MaxAttempts < 0 || r.MaxAttempts > 10 {
		return errors.New("max attempts out of range")
	}
	if r.Backoff.Base <= 0 || r.Backoff.Max < r.Backoff.Base {
		return errors.New("backoff invalid")
	}
	return nil
}
