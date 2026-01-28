package objects

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/IllumiKnowLabs/labstore/server/errs"
	"github.com/IllumiKnowLabs/labstore/server/helper"
)

// HeadObjectHandler: Head /:bucket/:key
func HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	key := r.PathValue("key")

	res, err := GetObject(bucket, key)
	if err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("head object", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			slog.Error("head object", "err", errNotFound)

			switch errNotFound.Type {
			case errs.ErrEntityTypeObject:
				errs.Handle(w, errs.S3NoSuchKey(errNotFound.Resource))
			case errs.ErrEntityTypeBucket:
				errs.Handle(w, errs.S3NoSuchBucket(bucket))
			default:
				errs.Handle(w, errs.S3InternalError())
			}

			return
		}

		slog.Error("head object", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}
	defer helper.CloseWithErr(res.Content, &err)

	buf := make([]byte, 512)

	n, err := res.Content.Read(buf)
	if err != nil && err != io.EOF {
		slog.Error("head object", "err", err)
		errs.Handle(w, errs.S3InternalError())
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
