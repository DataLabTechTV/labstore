package object

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

type GetObjectResult struct {
	Content      io.ReadSeekCloser
	ObjectSize   int
	DateModified time.Time
}

func GetObject(bucket, key string) (*GetObjectResult, error) {
	objPath := filepath.Join(config.Storage.ObjectsPath, bucket, key)

	file, err := os.Open(objPath)
	if err != nil {
		return nil, errs.S3NoSuchKey(key)
	}

	info, err := file.Stat()
	if err != nil {
		return nil, errs.S3InternalError("Couldn't compute file size")
	}

	res := &GetObjectResult{
		Content:      file,
		ObjectSize:   int(info.Size()),
		DateModified: info.ModTime(),
	}

	return res, nil
}

// GetObjectHandler: GET /:bucket/:key
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	res, err := GetObject(bucket, key)
	if err != nil {
		errs.Handle(w, err)
		return
	}
	defer helper.CloseWithErr(res.Content, &err)

	http.ServeContent(w, r, r.URL.Path, res.DateModified, res.Content)
}
