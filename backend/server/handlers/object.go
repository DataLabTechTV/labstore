package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/iam"
	"github.com/DataLabTechTV/labstore/backend/server"
)

// PutObject: PUT /:bucket/:key
func handlePutObject(
	w http.ResponseWriter,
	r *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "PutObject") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	bucketPath := filepath.Join(config.Env.StorageRoot, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		server.WriteS3Error(w, "NoSuchBucket", "Bucket does not exist", 404)
		return
	}
	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	os.MkdirAll(objDir, 0755)
	f, err := os.Create(objPath)
	if err != nil {
		server.WriteS3Error(w, "InternalError", "Failed to create object", 500)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		server.WriteS3Error(w, "InternalError", "Failed to write object", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetObject: GET /:bucket/:key
func handleGetObject(
	w http.ResponseWriter,
	r *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "GetObject") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)
	f, err := os.Open(objPath)
	if err != nil {
		server.WriteS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, key, time.Now(), f)
}

// DeleteObject: DELETE /:bucket/:key
func handleDeleteObject(
	w http.ResponseWriter,
	_ *http.Request,
	bucket,
	key,
	accessKey string,
) {
	if !iam.CheckPolicy(accessKey, bucket, "DeleteObject") {
		server.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)
	err := os.Remove(objPath)
	if err != nil {
		server.WriteS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
