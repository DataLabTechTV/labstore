package objects

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

func DeleteObject(bucket, key string) error {
	objPath := core.ObjectSystemPath(bucket, key)

	if !security.IsSubdir(config.Storage.ObjectsPath, objPath) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(bucket, key)}
	}

	err := os.Remove(objPath)
	if err != nil {
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(bucket, key)}
	}

	return nil
}

// DeleteObjectHandler: DELETE /:bucket/:key
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if err := DeleteObject(bucket, key); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("delete object", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			slog.Error("delete object", "err", errNotFound)
			errs.Handle(w, errs.S3NoSuchKey(errNotFound.Resource))
			return
		}

		slog.Error("delete object", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
