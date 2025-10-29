package service

import (
	"encoding/xml"
	"net/http"
	"os"
	"time"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
)

// !FIXME: move types to a proper location

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   struct {
		ID          string `xml:"ID"`
		DisplayName string `xml:"DisplayName"`
	}
	Buckets struct {
		Bucket []Bucket `xml:"Bucket"`
	}
}

// ListBuckets: GET /
func handleListBuckets(w http.ResponseWriter, _ *http.Request, accessKey string) {
	if !iam.CheckPolicy(accessKey, "", "ListBuckets") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}

	entries, err := os.ReadDir(config.Env.StorageRoot)
	if err != nil {
		core.WriteS3Error(w, "InternalError", "Failed to list buckets", 500)
		return
	}

	result := ListAllMyBucketsResult{}
	result.Owner.ID = accessKey
	result.Owner.DisplayName = accessKey

	for _, e := range entries {
		if e.IsDir() {
			b := Bucket{Name: e.Name(), CreationDate: time.Now().Format(time.RFC3339)}
			result.Buckets.Bucket = append(result.Buckets.Bucket, b)
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(result)
}
