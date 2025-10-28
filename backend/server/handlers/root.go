package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	"github.com/DataLabTechTV/labstore/backend/server"
	log "github.com/sirupsen/logrus"
)

func Start() {
	os.MkdirAll(config.Env.StorageRoot, 0755)
	http.HandleFunc("/", rootHandler)
	log.Infof("Starting minimal S3-compatible server on :%d", config.Env.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Env.Port), nil))
}

func NewRequestID() string {
	b := make([]byte, 8)
	helper.Must(rand.Read(b))
	return hex.EncodeToString(b)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	accessKey, err := server.VerifyAWSSigV4(r)
	if err != nil {
		server.WriteS3Error(w, "InvalidAccessKeyId", "Signature or access key invalid", 403)
		return
	}

	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)

	switch r.Method {
	case "PUT":
		if len(parts) == 1 {
			handlePutBucket(w, r, parts[0], accessKey)
		} else if len(parts) == 2 {
			handlePutObject(w, r, parts[0], parts[1], accessKey)
		}
	case "GET":
		if r.URL.Path == "/" || r.URL.Path == "" {
			handleListBuckets(w, r, accessKey)
		} else if len(parts) == 2 {
			handleGetObject(w, r, parts[0], parts[1], accessKey)
		} else {
			handleListObjects(w, r)
		}
	case "DELETE":
		if len(parts) == 1 {
			handleDeleteBucket(w, r, parts[0], accessKey)
		} else if len(parts) == 2 {
			handleDeleteObject(w, r, parts[0], parts[1], accessKey)
		}
	case "HEAD":
		if len(parts) == 2 {
			handleGetObject(w, r, parts[0], parts[1], accessKey)
		} else {
			handleListObjects(w, r)
		}
	default:
		server.WriteS3Error(w, "NotImplemented", "Operation not implemented", 501)
	}
}
