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

type BaseListBucketResult struct {
	XMLName        xml.Name `xml:"ListBucketResult"`
	Name           string
	Prefix         string
	MaxKeys        int
	Contents       []core.Object
	CommonPrefixes []CommonPrefixes
	IsTruncated    bool
}

type ListBucketResult struct {
	Marker     string
	NextMarker string
	BaseListBucketResult
}

type ListBucketResultV2 struct {
	KeyCount              int
	ContinuationToken     string
	NextContinuationToken string
	StartAfter            string
	BaseListBucketResult
}

type CommonPrefixes struct {
	Prefix string
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

func ListObjects(r *ListObjectsRequest) (*ListBucketResult, error) {
	slog.Debug("Processing ListObjects", "request", r)

	if !core.BucketExists(r.Bucket) {
		return nil, core.ErrorNoSuchBucket()
	}

	if r.Delimiter != "/" {
		return nil, errors.New("only '/' delimiters are supported by LabStore")
	}

	res := &ListBucketResult{
		BaseListBucketResult: BaseListBucketResult{
			Name:        r.Bucket,
			MaxKeys:     r.MaxKeys,
			IsTruncated: false,
		},
	}

	bucketPath := filepath.Join(config.Env.StorageRoot, r.Bucket)
	basePath := filepath.Join(bucketPath, r.Prefix)

	if !helper.FileExists(basePath) {
		return res, nil
	}

	err := res.list(bucketPath, basePath)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListObjectsV2(bucket string) (*ListBucketResult, error) {
	// TODO: implement
	slog.Debug("Processing ListObjectsV2")
	return &ListBucketResult{}, nil
}

// Lists objects as Contents, and directories as CommonPrefixes, for a given fs path
func (res *BaseListBucketResult) list(bucketPath, basePath string) error {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return errors.New("failed to list objects")
	}

	hash := md5.New()
	keyCount := 0

	for _, e := range entries {
		if e.IsDir() {
			path := filepath.Join(basePath, e.Name())
			key, err := filepath.Rel(bucketPath, path)
			if err != nil {
				return errors.New("could not resolve key")
			}
			key += "/"
			res.CommonPrefixes = append(res.CommonPrefixes, CommonPrefixes{Prefix: key})
			continue
		}

		name := e.Name()
		path := filepath.Join(basePath, name)
		key, err := filepath.Rel(bucketPath, path)
		if err != nil {
			return errors.New("could not resolve key")
		}

		info, err := e.Info()
		if err != nil {
			return fmt.Errorf("could not retrieve metadata: %s", key)
		}

		lastModified := core.Timestamp(info.ModTime())

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not read file: %s", key)
		}
		defer file.Close()

		if _, err := io.Copy(hash, file); err != nil {
			return fmt.Errorf("could not compute hash: %s", key)
		}
		eTag := hex.EncodeToString(hash.Sum(nil))

		size := info.Size()

		obj := core.Object{
			Key:          key,
			LastModified: lastModified,
			ETag:         eTag,
			Size:         size,
			// TODO add Owner when there is IAM (optional for V2)
		}

		res.Contents = append(res.Contents, obj)

		if keyCount++; keyCount >= res.MaxKeys {
			res.IsTruncated = true
			return nil
		}
	}

	return nil
}
