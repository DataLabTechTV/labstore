package core

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
)

func HandleError(w http.ResponseWriter, err error) {
	var s3Err *S3Error
	var iamErr *iam.IAMError

	if errors.As(err, &s3Err) {
		slog.Error("s3 error", "err", s3Err)
		WriteXML(w, s3Err.StatusCode, s3Err)
	} else if errors.As(err, &iamErr) {
		slog.Error("iam error", "err", iamErr)
		WriteXML(w, iamErr.StatusCode, iam.IAMErrorResponse{
			Error:     iamErr,
			RequestId: iamErr.RequestId,
		})
	} else {
		slog.Error("internal server error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
