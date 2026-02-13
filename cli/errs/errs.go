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

type ErrInsufficientArguments struct{}

func (e *ErrInsufficientArguments) Error() string {
	return "insufficient arguments"
}

type ErrNoProfileSelected struct{}

func (e *ErrNoProfileSelected) Error() string {
	return "no profile selected"
}

type ErrNoBucketSelected struct{}

func (e *ErrNoBucketSelected) Error() string {
	return "no bucket selected"
}

type ErrRemotePathNotSet struct{}

func (e *ErrRemotePathNotSet) Error() string {
	return "no remote path is set"
}

type ErrLocalPathNotSet struct{}

func (e *ErrLocalPathNotSet) Error() string {
	return "no local path is set"
}
