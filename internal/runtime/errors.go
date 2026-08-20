package runtime

import "fmt"

type Error struct {
	Code      string
	Operation string
	Cause     error
}

func (e Error) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Operation, e.Cause)
}
func (e Error) CodeValue() string { return e.Code }
func (e Error) Unwrap() error     { return e.Cause }
func Wrap(code, operation string, cause error) error {
	return Error{Code: code, Operation: operation, Cause: cause}
}
