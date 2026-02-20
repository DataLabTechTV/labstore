package errs

import "fmt"

type ErrHTTPStatusCode int

func (e ErrHTTPStatusCode) Error() string {
	return fmt.Sprintf("status code: %d", e)
}
