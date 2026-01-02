package bucket

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
)

func CreateBucket(bucket string) error {
	path := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.Storage.ObjectsPath, path) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if _, err := os.Stat(path); err == nil {
		return &errs.ErrExists{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}

// CreateBucket: PUT /:bucket
func PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")

	if err := CreateBucket(bucket); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("create bucket", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			slog.Error("create bucket", "err", errExists)
			errs.Handle(w, errs.S3BucketAlreadyExists())
			return
		}

		slog.Error("create bucket", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
