package object

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

func PutObject(bucket string, key string, reader io.Reader) error {
	bucketPath := filepath.Join(config.Storage.ObjectsPath, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return errs.S3NoSuchBucket(bucket)
	}

	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return errs.S3InternalError("Failed to create object directory")
	}

	f, err := os.Create(objPath)
	if err != nil {
		return errs.S3InternalError("Failed to create object")
	}
	defer helper.CloseWithErr(f, &err)

	writer := bufio.NewWriterSize(f, config.S3.IO.BufferSize)

	_, err = io.CopyBuffer(writer, reader, make([]byte, config.S3.IO.BufferSize))
	if err != nil {
		return errs.S3InternalError("Failed to write object")
	}

	return nil
}

// PutObjectHandler: PUT /:bucket/:key
func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if err := PutObject(bucket, key, r.Body); err != nil {
		errs.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
