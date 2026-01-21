package errs

import "fmt"

type RuntimeError struct{}

type ErrProfileNotFound struct {
	Name string
}

type ErrDefaultProfileNotSet struct{}

func (e *RuntimeError) Error() string {
	return "Runtime error"
}

func (e *ErrProfileNotFound) Error() string {
	return fmt.Sprintf("profile not found: %s", e.Name)
}

func (e *ErrDefaultProfileNotSet) Error() string {
	return "default profile not set"
}
