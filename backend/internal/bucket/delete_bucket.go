package bucket

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

func DeleteBucket(bucket string) error {
	path := filepath.Join(config.Storage.ObjectsPath, bucket)

	err := os.RemoveAll(path)
	if err != nil {
		return errs.S3NoSuchBucket(bucket)
	}

	return nil
}

// DeleteBucketHandler: DELETE /:bucket
func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")

	if err := DeleteBucket(bucket); err != nil {
		errs.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
