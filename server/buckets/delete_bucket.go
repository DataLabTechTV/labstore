package buckets

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/IllumiKnowLabs/labstore/server/core"
	"github.com/IllumiKnowLabs/labstore/server/errs"
	"github.com/IllumiKnowLabs/labstore/server/helper"
	"github.com/IllumiKnowLabs/labstore/server/security"
)

func DeleteBucket(bucket string) error {
	path := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.App.Server.Storage.ObjectsPath, path) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if !helper.FileExists(path) {
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeBucket, Resource: bucket}
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
