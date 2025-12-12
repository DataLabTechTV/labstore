package object

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

// HeadObjectHandler: Head /:bucket/:key
func HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	if !core.BucketExists(bucket) {
		errs.Handle(w, errs.S3NoSuchBucket(bucket))
		return
	}

	path := core.BucketKeyPath(bucket, key)
	if helper.IsDir(path) {
		errs.Handle(w, errs.S3NoSuchKey(key))
		return
	}

	res, err := GetObject(bucket, key)
	if err != nil {
		errs.Handle(w, err)
		return
	}
	defer helper.CloseWithErr(res.Content, &err)

	buf := make([]byte, 512)

	n, err := res.Content.Read(buf)
	if err != nil {
		errs.Handle(w, err)
		return
	}

	if _, err := res.Content.Seek(0, io.SeekStart); err != nil {
		slog.Warn("body seek", "err", err)
	}

	w.Header().Set("Content-Type", http.DetectContentType(buf[:n]))
	w.Header().Set("Content-Length", strconv.Itoa(res.ObjectSize))
	w.Header().Set("Last-Modified", res.DateModified.UTC().Format(http.TimeFormat))

	w.WriteHeader(http.StatusOK)
}
