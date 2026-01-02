package object

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func GetObject(bucket, key string) (*t.GetObjectResult, error) {
	bucketPath := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.Storage.ObjectsPath, bucketPath) {
		return nil, &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if !core.BucketExists(bucket) {
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	objPath := core.ObjectSystemPath(bucket, key)

	if !security.IsSubdir(config.Storage.ObjectsPath, objPath) {
		return nil, &errs.ErrForbidden{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(bucket, key)}
	}

	if helper.IsDir(objPath) {
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(bucket, key)}
	}

	file, err := os.Open(objPath)
	if err != nil {
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(bucket, key)}
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	res := &t.GetObjectResult{
		Content:      file,
		ObjectSize:   int(info.Size()),
		DateModified: info.ModTime(),
	}

	return res, nil
}

// GetObjectHandler: GET /:bucket/:key
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	res, err := GetObject(bucket, key)
	if err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("get object", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			slog.Error("get object", "err", errNotFound)

			switch errNotFound.Type {
			case errs.ErrEntityTypeObject:
				errs.Handle(w, errs.S3NoSuchKey(errNotFound.Resource))
			case errs.ErrEntityTypeBucket:
				errs.Handle(w, errs.S3NoSuchBucket(bucket))
			default:
				errs.Handle(w, errs.S3InternalError())
			}

			return
		}

		slog.Error("get object", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}
	defer helper.CloseWithErr(res.Content, &err)

	http.ServeContent(w, r, r.URL.Path, res.DateModified, res.Content)
}
