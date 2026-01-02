package bucket

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
)

func DeleteBucket(bucket string) error {
	path := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.Storage.ObjectsPath, path) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	err := os.RemoveAll(path)
	if err != nil {
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	return nil
}

// DeleteBucketHandler: DELETE /:bucket
func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")

	if err := DeleteBucket(bucket); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("delete bucket", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			slog.Error("delete bucket", "err", errNotFound)
			errs.Handle(w, errs.S3NoSuchBucket(errNotFound.Resource))
			return
		}

		slog.Error("delete bucket", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
