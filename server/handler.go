package server

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Received probe: ", r.URL.Path)

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Server", "LabStore")
	w.Header().Set("x-amz-request-id", "12345")

	w.WriteHeader(http.StatusOK)
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

// Create bucket: PUT /bucket
func handlePutBucket(w http.ResponseWriter, _ *http.Request, bucket string, accessKey string) {
	if !checkPolicy(accessKey, bucket, "CreateBucket") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(storageRoot, bucket)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			writeS3Error(w, "InternalError", "Could not create bucket", 500)
			return
		}
	}
	w.WriteHeader(200)
}

// Delete bucket: DELETE /bucket
func handleDeleteBucket(w http.ResponseWriter, _ *http.Request, bucket string, accessKey string) {
	if !checkPolicy(accessKey, bucket, "DeleteBucket") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	path := filepath.Join(storageRoot, bucket)
	err := os.Remove(path)
	if err != nil {
		writeS3Error(w, "NoSuchBucket", "Bucket does not exist or not empty", 404)
		return
	}
	w.WriteHeader(204)
}

// Upload object: PUT /bucket/key
func handlePutObject(w http.ResponseWriter, r *http.Request, bucket, key, accessKey string) {
	if !checkPolicy(accessKey, bucket, "PutObject") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	bucketPath := filepath.Join(storageRoot, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		writeS3Error(w, "NoSuchBucket", "Bucket does not exist", 404)
		return
	}
	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	os.MkdirAll(objDir, 0755)
	f, err := os.Create(objPath)
	if err != nil {
		writeS3Error(w, "InternalError", "Failed to create object", 500)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		writeS3Error(w, "InternalError", "Failed to write object", 500)
		return
	}
	w.WriteHeader(200)
}

// Get object: GET /bucket/key
func handleGetObject(w http.ResponseWriter, r *http.Request, bucket, key, accessKey string) {
	if !checkPolicy(accessKey, bucket, "GetObject") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(storageRoot, bucket, key)
	f, err := os.Open(objPath)
	if err != nil {
		writeS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, key, time.Now(), f)
}

// Delete object: DELETE /bucket/key
func handleDeleteObject(w http.ResponseWriter, _ *http.Request, bucket, key, accessKey string) {
	if !checkPolicy(accessKey, bucket, "DeleteObject") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	objPath := filepath.Join(storageRoot, bucket, key)
	err := os.Remove(objPath)
	if err != nil {
		writeS3Error(w, "NoSuchKey", "Object not found", 404)
		return
	}
	w.WriteHeader(204)
}

// List buckets: GET /
func handleListBuckets(w http.ResponseWriter, _ *http.Request, accessKey string) {
	if !checkPolicy(accessKey, "", "ListBuckets") {
		writeS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}
	// For MVP: just list bucket dirs in storageRoot
	entries, err := os.ReadDir(storageRoot)
	if err != nil {
		writeS3Error(w, "InternalError", "Failed to list buckets", 500)
		return
	}
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

// --- Main router ---

func handler(w http.ResponseWriter, r *http.Request) {
	accessKey, ok := verifyAWSSigV4(r)
	if !ok {
		writeS3Error(w, "InvalidAccessKeyId", "Signature or access key invalid", 403)
		return
	}

	// Parse URL path: /bucket or /bucket/key
	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)

	switch r.Method {
	case "PUT":
		if len(parts) == 1 {
			// PUT Bucket
			handlePutBucket(w, r, parts[0], accessKey)
			return
		} else if len(parts) == 2 {
			// PUT Object
			handlePutObject(w, r, parts[0], parts[1], accessKey)
			return
		}
	case "GET":
		if r.URL.Path == "/" || r.URL.Path == "" {
			handleListBuckets(w, r, accessKey)
			return
		} else if len(parts) == 2 {
			handleGetObject(w, r, parts[0], parts[1], accessKey)
			return
		} else {
			handleRoot(w, r)
			return
		}
	case "DELETE":
		if len(parts) == 1 {
			handleDeleteBucket(w, r, parts[0], accessKey)
			return
		} else if len(parts) == 2 {
			handleDeleteObject(w, r, parts[0], parts[1], accessKey)
			return
		}
	case "HEAD":
		handleRoot(w, r)
		return
	}

	writeS3Error(w, "NotImplemented", "Operation not implemented", 501)
}

func Start() {
	os.MkdirAll(storageRoot, 0755)
	http.HandleFunc("/", handler)
	log.Info("Starting minimal S3-compatible server on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
