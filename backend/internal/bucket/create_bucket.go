package bucket

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

func CreateBucket(bucket string) error {
	path := filepath.Join(config.Storage.ObjectsPath, bucket)

	if _, err := os.Stat(path); err == nil {
		return errs.S3BucketAlreadyExists()
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("could not create bucket: %w", err)
	}

	return nil
}

// CreateBucket: PUT /:bucket
func PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")

	if err := CreateBucket(bucket); err != nil {
		errs.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
