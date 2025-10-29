package bucket

import (
	"net/http"

	"github.com/DataLabTechTV/labstore/backend/middleware"
	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
)

// ListObjects: GET /:bucket
func ListObjects(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.NewRequestID()
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
