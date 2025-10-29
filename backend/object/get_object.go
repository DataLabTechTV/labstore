package object

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/iam"
)

// GetObject: GET /:bucket/:key
func GetObject(
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
	f, err := os.Open(objPath)
	if err != nil {
		core.WriteS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, key, time.Now(), f)
}
