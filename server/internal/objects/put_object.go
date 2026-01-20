package objects

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/security"
)

func PutObject(bucket string, key string, reader io.Reader) error {
	bucketPath := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.App.Server.Storage.ObjectsPath, bucketPath) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeBucket, Resource: bucket}
	}

	objPath := filepath.Join(bucketPath, key)

	if helper.IsDir(objPath) {
		return nil
	}

	objDir := filepath.Dir(objPath)
	if helper.FileExists(objDir) && !helper.IsDir(objDir) {
		return &errs.ErrNotDirectory{Path: objDir}
	}

	if err := os.MkdirAll(objDir, 0o755); err != nil {
		return err
	}

	file, err := os.Create(objPath)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(file, &err)

	writer := bufio.NewWriterSize(file, config.App.Server.S3.IO.BufferSize)
	buf := make([]byte, config.App.Server.S3.IO.BufferSize)

	_, err = io.CopyBuffer(writer, reader, buf)
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
			slog.Error("put object", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			if errNotFound.Type == errs.ErrEntityTypeBucket {
				slog.Error("put object", "err", errNotFound)
				errs.Handle(w, errs.S3NoSuchBucket(errNotFound.Resource))
				return
			}
		}

		var errNotDirectory *errs.ErrNotDirectory
		if errors.As(err, &errNotDirectory) {
			slog.Error("put object", "err", errNotDirectory)
			errs.Handle(w, errs.S3OperationAborted(bucket, key))
			return
		}

		slog.Error("put object", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
