package object

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

func DeleteObject(bucket, key string) error {
	objPath := filepath.Join(config.Storage.ObjectsPath, bucket, key)

	err := os.Remove(objPath)
	if err != nil {
		return errs.S3NoSuchKey(key)
	}

	return nil
}

// DeleteObjectHandler: DELETE /:bucket/:key
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if err := DeleteObject(bucket, key); err != nil {
		errs.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
