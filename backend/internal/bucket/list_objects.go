package bucket

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
)

const DefaultMaxKeys = 1000

type ListObjectsRequest struct {
	Bucket    string
	Prefix    string
	Delimiter string
	MaxKeys   int
	// TODO: ...
}

type ListBucketResult struct {
	XMLName        xml.Name `xml:"ListBucketResult"`
	Name           string
	Prefix         string
	Marker         string
	MaxKeys        int
	IsTruncated    bool
	Contents       []core.Object
	CommonPrefixes []CommonPrefixes
}

type CommonPrefixes struct {
	Prefix string
}

func ListObjects(r *ListObjectsRequest) (*ListBucketResult, error) {
	slog.Debug("Processing ListObjects", "request", r)

	if !core.BucketExists(r.Bucket) {
		return nil, core.ErrorNoSuchBucket()
	}

	if r.Delimiter != "/" {
		return nil, errors.New("only '/' delimiters are supported by LabStore")
	}

	res := &ListBucketResult{
		Name:        r.Bucket,
		MaxKeys:     r.MaxKeys,
		Contents:    []core.Object{},
		IsTruncated: false,
	}

	hash := md5.New()
	keyCount := 0
	bucketPath := filepath.Join(config.Env.StorageRoot, r.Bucket)
	basePath := filepath.Join(bucketPath, r.Prefix)

	if !helper.FileExists(basePath) {
		return res, nil
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, errors.New("failed to list objects")
	}

	for _, e := range entries {
		if e.IsDir() {
			path := filepath.Join(basePath, e.Name())
			key, err := filepath.Rel(bucketPath, path)
			if err != nil {
				return nil, errors.New("could not resolve key")
			}
			key += "/"
			res.CommonPrefixes = append(res.CommonPrefixes, CommonPrefixes{Prefix: key})
			continue
		}

		name := e.Name()
		path := filepath.Join(basePath, name)
		key, err := filepath.Rel(bucketPath, path)
		if err != nil {
			return nil, errors.New("could not resolve key")
		}

		info, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve metadata: %s", key)
		}

		lastModified := core.Timestamp(info.ModTime())

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("could not read file: %s", key)
		}
		defer file.Close()

		if _, err := io.Copy(hash, file); err != nil {
			return nil, fmt.Errorf("could not compute hash: %s", key)
		}
		eTag := hex.EncodeToString(hash.Sum(nil))

		size := info.Size()

		obj := core.Object{
			Key:          key,
			LastModified: lastModified,
			ETag:         eTag,
			Size:         size,
			// TODO add missing Owner when there is proper IAM
		}

		res.Contents = append(res.Contents, obj)

		if keyCount++; keyCount >= res.MaxKeys {
			res.IsTruncated = true
			return res, nil
		}
	}

	return res, nil
}

func ListObjectsV2(bucket string) (*ListBucketResult, error) {
	// TODO: implement
	slog.Debug("Processing ListObjectsV2")
	return &ListBucketResult{}, nil
}

// ListObjectsHandler: GET /:bucket
func ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	requestID := middleware.NewRequestID()

	q := r.URL.Query()

	prefix := q.Get("prefix")
	delimiter := q.Get("delimiter")

	var err error
	var maxKeys int
	var res *ListBucketResult

	if mk := q.Get("maxKeys"); mk == "" {
		maxKeys = DefaultMaxKeys
	} else {
		if maxKeys, err = strconv.Atoi(mk); err != nil {
			slog.Warn("Invalid max-keys value, using default...")
			maxKeys = DefaultMaxKeys
		}
	}

	if q.Get("list-type") == "2" {
		res, err = ListObjectsV2(bucket)
	} else {
		r := &ListObjectsRequest{
			Bucket:    bucket,
			Prefix:    prefix,
			Delimiter: delimiter,
			MaxKeys:   maxKeys,
		}
		res, err = ListObjects(r)
	}

	if err != nil {
		core.HandleError(w, err)
		return
	}

	w.Header().Set("Server", "LabStore")
	w.Header().Set("X-Amz-Request-Id", requestID)

	core.WriteXML(w, http.StatusOK, res)
}
