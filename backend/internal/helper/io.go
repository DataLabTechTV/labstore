package helper

import (
	"io"
	"log/slog"
)

func CloseWithErr(closer io.Closer, prevErr *error) {
	if closeErr := closer.Close(); closeErr != nil {
		if *prevErr == nil {
			*prevErr = closeErr
		} else {
			slog.Warn("close failed", "err", closeErr)
		}
	}
}
