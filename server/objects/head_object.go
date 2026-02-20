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
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			slog.Error("head object", "err", errNotFound)

			switch errNotFound.Type {
			case errs.ErrEntityTypeObject, errs.ErrEntityTypeBucket:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		slog.Error("head object", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer helper.CloseWithErr(res.Content, &err)

	buf := make([]byte, 512)

	n, err := res.Content.Read(buf)
	if err != nil && err != io.EOF {
		slog.Error("head object", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
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
