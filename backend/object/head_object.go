package object

import (
	"net/http"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/iam"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
)

// HeadObject: Head /:bucket/:key
func HeadObject(
	w http.ResponseWriter,
	r *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "GetObject") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)
	if helper.FileExists(objPath) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
