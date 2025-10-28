package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/iam"
	"github.com/DataLabTechTV/labstore/backend/server"
	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
)

// CreateBucket: PUT /:bucket
func handlePutBucket(
	w http.ResponseWriter,
	_ *http.Request,
	bucket string,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "CreateBucket") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(config.Env.StorageRoot, bucket)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			server.WriteS3Error(w, "InternalError", "Could not create bucket", 500)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// ListObjects: GET /:bucket
func handleListObjects(w http.ResponseWriter, r *http.Request) {
	requestID := NewRequestID()
	log.Info("Received probe: ", r.URL.Path)

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Server", "LabStore")
	w.Header().Set("X-Amz-Request-Id", requestID)

	w.WriteHeader(http.StatusNotFound)

	// !FIXME: ListAllMyBucketsResult is for ListBuckets, not ListObjects
	w.Write([]byte(
		heredoc.Doc(`
			<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
			  <Owner>
			    <ID>admin</ID>
			    <DisplayName>admin</DisplayName>
			  </Owner>
			  <Buckets>
			  </Buckets>
			</ListAllMyBucketsResult>
		`),
	))
}

// DeleteBucket: DELETE /:bucket
func handleDeleteBucket(
	w http.ResponseWriter,
	_ *http.Request,
	bucket string,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "DeleteBucket") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(config.Env.StorageRoot, bucket)
	err := os.Remove(path)
	if err != nil {
		server.WriteS3Error(w, "NoSuchBucket", "Bucket does not exist or not empty", 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
