package runtime

import "errors"

var ErrTimeout = errors.New("runtime timeout")
var ErrCancelled = errors.New("runtime cancelled")
var ErrMemory = errors.New("runtime memory limit")
var ErrInstructions = errors.New("runtime instruction limit")

func Retryable(err error) bool {
	return errors.Is(err, ErrTimeout) || errors.Is(err, ErrCancelled)
}
func Classify(err error) string {
	if Retryable(err) {
		return "retryable"
	}
	return "terminal"
}
