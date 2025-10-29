package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/auth"
	"github.com/DataLabTechTV/labstore/backend/bucket"
	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/core"
	"github.com/DataLabTechTV/labstore/backend/object"
	log "github.com/sirupsen/logrus"
)

func Start() {
	os.MkdirAll(config.Env.StorageRoot, 0755)
	http.HandleFunc("/", rootHandler)
	log.Infof("Starting minimal S3-compatible server on :%d", config.Env.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Env.Port), nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	accessKey, err := auth.VerifyAWSSigV4(r)
	if err != nil {
		core.WriteS3Error(w, "InvalidAccessKeyId", "Signature or access key invalid", 403)
		return
	}

	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)

	switch r.Method {
	case "PUT":
		if len(parts) == 1 {
			bucket.PutBucket(w, r, parts[0], accessKey)
		} else if len(parts) == 2 {
			object.PutObject(w, r, parts[0], parts[1], accessKey)
		}
	case "GET":
		if r.URL.Path == "/" || r.URL.Path == "" {
			handleListBuckets(w, r, accessKey)
		} else if len(parts) == 2 {
			object.GetObject(w, r, parts[0], parts[1], accessKey)
		} else {
			bucket.ListObjects(w, r)
		}
	case "DELETE":
		if len(parts) == 1 {
			bucket.DeleteBucket(w, r, parts[0], accessKey)
		} else if len(parts) == 2 {
			object.DeleteObject(w, r, parts[0], parts[1], accessKey)
		}
	case "HEAD":
		if r.URL.Path == "/" || r.URL.Path == "" {
			bucket.HeadBucket(w, r, accessKey)
		} else if len(parts) == 2 {
			object.HeadObject(w, r, parts[0], parts[1], accessKey)
		}
	default:
		core.WriteS3Error(w, "NotImplemented", "Operation not implemented", 501)
	}
}
