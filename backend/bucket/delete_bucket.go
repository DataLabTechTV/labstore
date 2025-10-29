package bucket

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/iam"
)

// DeleteBucket: DELETE /:bucket
func DeleteBucket(
	w http.ResponseWriter,
	_ *http.Request,
	bucket string,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "DeleteBucket") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(config.Env.StorageRoot, bucket)
	err := os.Remove(path)
	if err != nil {
		core.WriteS3Error(w, "NoSuchBucket", "Bucket does not exist or not empty", 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
