package object

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
)

// DeleteObject: DELETE /:bucket/:key
func DeleteObject(
	w http.ResponseWriter,
	_ *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "DeleteObject") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)
	err := os.Remove(objPath)
	if err != nil {
		core.WriteS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
