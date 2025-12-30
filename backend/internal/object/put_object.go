package object

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

func PutObject(bucket string, key string, reader io.Reader) error {
	bucketPath := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.Storage.ObjectsPath, bucketPath) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return err
	}

	f, err := os.Create(objPath)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(f, &err)

	writer := bufio.NewWriterSize(f, config.S3.IO.BufferSize)

	_, err = io.CopyBuffer(writer, reader, make([]byte, config.S3.IO.BufferSize))
	if err != nil {
		return err
	}

	return nil
}

// PutObjectHandler: PUT /:bucket/:key
func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if err := PutObject(bucket, key, r.Body); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("put object: unsafe path", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			if errNotFound.Type == errs.ErrEntityTypeBucket {
				slog.Error("put object: bucket not found", "err", errNotFound)
				errs.Handle(w, errs.S3NoSuchBucket(errNotFound.Resource))
				return
			}
		}

		slog.Error("put object", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
