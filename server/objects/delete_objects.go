package objects

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/IllumiKnowLabs/labstore/server/core"
	"github.com/IllumiKnowLabs/labstore/server/errs"
	"github.com/IllumiKnowLabs/labstore/server/helper"
	"github.com/IllumiKnowLabs/labstore/server/security"
	"github.com/IllumiKnowLabs/labstore/server/types"
)

func DeleteObjects(bucket string, r *types.DeleteObjectsRequest) *types.DeleteResult {
	res := &types.DeleteResult{}
	bucketPath := core.BucketSystemPath(bucket)

	for _, obj := range r.Object {
		objPath := filepath.Join(bucketPath, obj.Key)

		if !security.IsSubdir(config.App.Server.Storage.ObjectsPath, objPath) {
			res.Error = append(res.Error, *errs.S3AccessDenied())
		}

		if !helper.FileExists(objPath) {
			res.Error = append(res.Error, *errs.S3NoSuchKey(obj.Key))
			continue
		}

		err := os.RemoveAll(objPath)
		if err != nil {
			res.Error = append(res.Error, *errs.S3NoSuchKey(obj.Key))
			continue
		}

		deleted := types.DeletedObject{
			DeleteMarker: false,
			Key:          obj.Key,
		}
		res.Deleted = append(res.Deleted, deleted)
	}

	return res
}

// DeleteObjectsHandler: POST /:bucket?delete=
func DeleteObjectsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if !q.Has("delete") {
		http.Error(w, "Only delete requests are supported", http.StatusBadRequest)
		return
	}

	bucket := r.PathValue("bucket")

	var req types.DeleteObjectsRequest
	err := helper.ReadXML(r.Body, &req)
	if err != nil {
		slog.Error("delete objects", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	slog.Debug("delete objects", "request", req)

	resp := DeleteObjects(bucket, &req)
	helper.WriteXMLResponse(w, http.StatusOK, resp)
}
