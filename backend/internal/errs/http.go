package errs

import "fmt"

func HTTPMissingQueryParam(param string) error {
	return fmt.Errorf("missing query parameter: %s", param)
}
