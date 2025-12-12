package errs

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
)

func Handle(w http.ResponseWriter, err error) {
	var s3Err *S3Error
	var iamErr *IAMError

	if errors.As(err, &s3Err) {
		slog.Error("s3 error", "err", s3Err)
		core.WriteXML(w, s3Err.StatusCode, s3Err)
	} else if errors.As(err, &iamErr) {
		slog.Error("iam error", "err", iamErr)
		core.WriteXML(w, iamErr.StatusCode, IAMErrorResponse{
			Error:     iamErr,
			RequestId: iamErr.RequestId,
		})
	} else {
		slog.Error("internal server error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
