package service

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"

	"github.com/djherbis/times"
)

func ListBuckets(accessKey string) (*types.ListAllMyBucketsResult, error) {
	entries, err := os.ReadDir(config.App.Server.Storage.ObjectsPath)
	if err != nil {
		return nil, err
	}

	res := types.ListAllMyBucketsResult{}
	res.Owner.ID = accessKey
	res.Owner.DisplayName = accessKey

	for _, entry := range entries {
		if entry.IsDir() {
			var birthDate time.Time

			path := filepath.Join(config.App.Server.Storage.ObjectsPath, entry.Name())

			stat, err := times.Stat(path)
			if err != nil {
				birthDate = time.Unix(0, 0).UTC()
			} else {
				birthDate = stat.BirthTime()
			}

			bucket := types.Bucket{
				Name:         entry.Name(),
				CreationDate: types.Timestamp(birthDate),
			}
			res.Buckets.Bucket = append(res.Buckets.Bucket, bucket)
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
