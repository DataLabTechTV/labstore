package bucket

import (
	"encoding/xml"
	"net/http"

	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
)

type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	IsTruncated bool
	Marker      string
	Name        string
	Prefix      string
	MaxKeys     int
}

func ListBucket(bucket string) (*ListBucketResult, error) {
	// TODO: implement
	return &ListBucketResult{}, nil
}

// ListObjectsHandler: GET /:bucket
func ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	requestID := middleware.NewRequestID()

	res, err := ListBucket(bucket)
	if err != nil {
		core.HandleError(w, err)
		return
	}

	w.Header().Set("Server", "LabStore")
	w.Header().Set("X-Amz-Request-Id", requestID)

	core.WriteXML(w, http.StatusOK, res)
}
