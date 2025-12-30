package service

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
)

// !FIXME: move types to a proper location

type Bucket struct {
	Name         string
	CreationDate string
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   struct {
		ID          string
		DisplayName string
	}
	Buckets struct {
		Bucket []Bucket
	}
}

func ListBuckets(accessKey string) (*ListAllMyBucketsResult, error) {
	entries, err := os.ReadDir(config.Storage.ObjectsPath)
	if err != nil {
		return nil, err
	}

	res := ListAllMyBucketsResult{}
	res.Owner.ID = accessKey
	res.Owner.DisplayName = accessKey

	for _, e := range entries {
		if e.IsDir() {
			b := Bucket{Name: e.Name(), CreationDate: time.Now().Format(time.RFC3339)}
			res.Buckets.Bucket = append(res.Buckets.Bucket, b)
		}
	}

	return &res, nil
}

// ListBuckets: GET /
func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	accessKey := middleware.GetRequestAccessKey(r)

	res, err := ListBuckets(accessKey)
	if err != nil {
		slog.Error("list buckets", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	core.WriteXML(w, http.StatusOK, res)
}
