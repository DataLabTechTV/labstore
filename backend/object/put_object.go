package object

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/iam"
)

// PutObject: PUT /:bucket/:key
func PutObject(
	w http.ResponseWriter,
	r *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "PutObject") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	bucketPath := filepath.Join(config.Env.StorageRoot, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		core.WriteS3Error(w, "NoSuchBucket", "Bucket does not exist", 404)
		return
	}
	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	os.MkdirAll(objDir, 0755)
	f, err := os.Create(objPath)
	if err != nil {
		core.WriteS3Error(w, "InternalError", "Failed to create object", 500)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		core.WriteS3Error(w, "InternalError", "Failed to write object", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}
