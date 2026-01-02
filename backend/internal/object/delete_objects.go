package object

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
)

type DeleteObjectsRequest struct {
	XMLName xml.Name `xml:"Delete"`
	Object  []core.ObjectIdentifier
	Quiet   bool
}

type DeleteResult struct {
	Deleted []DeletedObject
	Error   []errs.S3Error
}

type DeletedObject struct {
	DeleteMarker          bool
	DeleteMarkerVersionId string
	Key                   string
	VersionId             string
}

func (req DeleteObjectsRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("XMLName", req.XMLName),
		slog.Int("Objects", len(req.Object)),
		slog.Bool("Quiet", req.Quiet),
	)
}

func DeleteObjects(bucket string, r *DeleteObjectsRequest) *DeleteResult {
	res := &DeleteResult{}
	bucketPath := core.BucketSystemPath(bucket)

	for _, obj := range r.Object {
		objPath := filepath.Join(bucketPath, obj.Key)

		if !security.IsSubdir(config.Storage.ObjectsPath, objPath) {
			res.Error = append(res.Error, *errs.S3AccessDenied())
		}

		err := os.RemoveAll(objPath)
		if err != nil {
			res.Error = append(res.Error, *errs.S3NoSuchKey(obj.Key))
			continue
		}

		deleted := DeletedObject{
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

	var req DeleteObjectsRequest
	err := core.ReadXML(w, r, &req)
	if err != nil {
		slog.Error("delete objects", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	slog.Debug("delete objects", "request", req)

	resp := DeleteObjects(bucket, &req)
	core.WriteXML(w, http.StatusOK, resp)
}
