package service

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func ListBuckets(accessKey string) (*t.ListAllMyBucketsResult, error) {
	entries, err := os.ReadDir(config.Storage.ObjectsPath)
	if err != nil {
		return nil, err
	}

	res := t.ListAllMyBucketsResult{}
	res.Owner.ID = accessKey
	res.Owner.DisplayName = accessKey

	for _, e := range entries {
		if e.IsDir() {
			b := t.Bucket{Name: e.Name(), CreationDate: time.Now().Format(time.RFC3339)}
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

	helper.WriteXMLResponse(w, http.StatusOK, res)
}
