package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// --- Basic config and types ---

var (
	storageRoot = "./data" // local dir to store buckets & objects

	// In-memory user with access key & secret key (hardcoded for MVP)
	users = map[string]string{
		"myaccesskey": "mysecretkey",
	}

	// Simple user policies (for MVP: allow all operations)
	userPolicies = map[string]func(string, string) bool{
		"myaccesskey": func(bucket, op string) bool {
			return true
		},
	}
)

// S3 XML Error response
type S3Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string
	Message   string
	RequestId string
	HostId    string
}

func writeS3Error(w http.ResponseWriter, code, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errResp := S3Error{
		Code:    code,
		Message: message,
	}
	xml.NewEncoder(w).Encode(errResp)
}

// --- AWS Signature V4 Verification ---

// For simplicity, we only check Authorization header with access key and sign string presence.
// Full signature verification is complex, so this is a stub.

func verifyAWSSigV4(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	// Example: Authorization: AWS4-HMAC-SHA256 Credential=myaccesskey/20230629/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-date, Signature=...
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return "", false
	}
	// Parse Credential= part
	parts := strings.Split(auth, ",")
	var credential string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if after, ok := strings.CutPrefix(p, "Credential="); ok {
			credential = after
			break
		}
	}
	if credential == "" {
		return "", false
	}
	accessKey := strings.Split(credential, "/")[0]
	_, ok := users[accessKey]
	if !ok {
		return "", false
	}
	// Here you'd do full signature verification with secretKey, canonical request, etc.
	// For MVP, we just check presence of accessKey in users.
	return accessKey, true
}

// --- Simple IAM policy enforcement (allow all for MVP) ---

func checkPolicy(accessKey, bucket, op string) bool {
	if polFunc, ok := userPolicies[accessKey]; ok {
		return polFunc(bucket, op)
	}
	return false
}

// --- Handlers ---

// Create bucket: PUT /bucket
func handlePutBucket(w http.ResponseWriter, r *http.Request, bucket string, accessKey string) {
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
func handleDeleteBucket(w http.ResponseWriter, r *http.Request, bucket string, accessKey string) {
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
func handleDeleteObject(w http.ResponseWriter, r *http.Request, bucket, key, accessKey string) {
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
func handleListBuckets(w http.ResponseWriter, r *http.Request, accessKey string) {
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
		}
	case "DELETE":
		if len(parts) == 1 {
			handleDeleteBucket(w, r, parts[0], accessKey)
			return
		} else if len(parts) == 2 {
			handleDeleteObject(w, r, parts[0], parts[1], accessKey)
			return
		}
	}

	writeS3Error(w, "NotImplemented", "Operation not implemented", 501)
}

func main() {
	log.Println("Welcome to Lab Store, by https://youtube.com/@DataLabTechTV")
	os.MkdirAll(storageRoot, 0755)
	http.HandleFunc("/", handler)
	log.Println("Starting minimal S3-compatible server on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
