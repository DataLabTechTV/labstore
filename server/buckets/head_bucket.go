package buckets

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/IllumiKnowLabs/labstore/server/core"
	"github.com/IllumiKnowLabs/labstore/server/security"
)

func HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")

	bucketPath := core.BucketSystemPath(bucket)

	if !security.IsSubdir(config.App.Server.Storage.ObjectsPath, bucketPath) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !core.BucketExists(bucket) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
