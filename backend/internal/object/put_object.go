package object

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

func PutObject(bucket string, key string, reader io.Reader) error {
	bucketPath := filepath.Join(config.S3.Storage.Path, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return core.ErrNoSuchBucket(bucket)
	}

	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return core.ErrInternalError("Failed to create object directory")
	}

	f, err := os.Create(objPath)
	if err != nil {
		return core.ErrInternalError("Failed to create object")
	}
	defer helper.CloseWithErr(f, &err)

	writer := bufio.NewWriterSize(f, config.S3.Perf.BufferSize)

	_, err = io.CopyBuffer(writer, reader, make([]byte, config.S3.Perf.BufferSize))
	if err != nil {
		return core.ErrInternalError("Failed to write object")
	}

	return nil
}

// PutObjectHandler: PUT /:bucket/:key
func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if err := PutObject(bucket, key, r.Body); err != nil {
		core.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
