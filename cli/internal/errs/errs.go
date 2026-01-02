package errs

import "fmt"

type ErrProfileNotFound struct {
	Name string
}

type ErrDefaultProfileNotSet struct{}

func (e *ErrProfileNotFound) Error() string {
	return fmt.Sprintf("profile not found: %s", e.Name)
}

func (e *ErrDefaultProfileNotSet) Error() string {
	return "default profile not set"
}
