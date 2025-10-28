package handlers

import (
	"encoding/xml"
	"net/http"
	"os"
	"time"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/iam"
	"github.com/DataLabTechTV/labstore/backend/server"
)

// ListBuckets: GET /
func handleListBuckets(w http.ResponseWriter, _ *http.Request, accessKey string) {
	if !iam.CheckPolicy(accessKey, "", "ListBuckets") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}

	entries, err := os.ReadDir(config.Env.StorageRoot)
	if err != nil {
		server.WriteS3Error(w, "InternalError", "Failed to list buckets", 500)
		return
	}

	// !FIXME: move types out of handler
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
