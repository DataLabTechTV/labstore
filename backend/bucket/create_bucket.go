package bucket

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/iam"
)

// CreateBucket: PUT /:bucket
func PutBucket(
	w http.ResponseWriter,
	_ *http.Request,
	bucket string,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "CreateBucket") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(config.Env.StorageRoot, bucket)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			core.WriteS3Error(w, "InternalError", "Could not create bucket", 500)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
