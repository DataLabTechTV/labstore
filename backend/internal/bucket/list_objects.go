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
	"strings"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
)

const DefaultMaxKeys = 1000
const DefaultDelimiter = "/"

type BaseListObjectsRequest struct {
	Bucket    string
	Prefix    string
	Delimiter string
	MaxKeys   int
}

type ListObjectsRequest struct {
	BaseListObjectsRequest
	Marker string
}

type ListObjectsRequestV2 struct {
	BaseListObjectsRequest
	ContinuationToken string
	StartAfter        string
	FetchOwner        bool
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
	BaseListBucketResult
	Marker     string
	NextMarker string
}

type ListBucketResultV2 struct {
	BaseListBucketResult
	KeyCount              int
	ContinuationToken     string
	NextContinuationToken string
	StartAfter            string
}

type CommonPrefixes struct {
	Prefix string
}

// ListObjectsHandler: GET /:bucket
func ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	var res any
	var err error
	var delimiter string
	var maxKeys int

	bucket := r.PathValue("bucket")
	requestID := middleware.NewRequestID()

	q := r.URL.Query()

	prefix := q.Get("prefix")

	if d := q.Get("delimiter"); d == "" {
		delimiter = DefaultDelimiter
	} else {
		delimiter = d
	}

	if mk := q.Get("maxKeys"); mk == "" {
		maxKeys = DefaultMaxKeys
	} else {
		if maxKeys, err = strconv.Atoi(mk); err != nil {
			slog.Warn("Invalid max-keys value, using default...")
			maxKeys = DefaultMaxKeys
		}
	}

	rBase := BaseListObjectsRequest{
		Bucket:    bucket,
		Prefix:    prefix,
		Delimiter: delimiter,
		MaxKeys:   maxKeys,
	}

	if q.Get("list-type") == "2" {
		r := &ListObjectsRequestV2{
			BaseListObjectsRequest: rBase,
		}
		res, err = ListObjectsV2(r)

	} else {
		r := &ListObjectsRequest{
			BaseListObjectsRequest: rBase,
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

	err := res.list(bucketPath, r.Prefix, r.Delimiter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListObjectsV2(r *ListObjectsRequestV2) (*ListBucketResultV2, error) {
	slog.Debug("Processing ListObjectsV2", "request", r)

	if !core.BucketExists(r.Bucket) {
		return nil, core.ErrorNoSuchBucket()
	}

	if r.Delimiter != "/" {
		return nil, errors.New("only '/' delimiters are supported by LabStore")
	}

	res := &ListBucketResultV2{
		BaseListBucketResult: BaseListBucketResult{
			Name:        r.Bucket,
			MaxKeys:     r.MaxKeys,
			IsTruncated: false,
		},
	}

	bucketPath := filepath.Join(config.Env.StorageRoot, r.Bucket)

	err := res.list(bucketPath, r.Prefix, r.Delimiter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Lists objects as Contents, and directories as CommonPrefixes, for a given fs path
func (res *BaseListBucketResult) list(bucketPath, prefix, delimiter string) error {
	var paths []string
	var basePath string

	if strings.HasSuffix(prefix, delimiter) {
		// full prefix
		var entries []os.DirEntry
		basePath = filepath.Join(bucketPath, prefix)
		entries, err := os.ReadDir(basePath)

		if err != nil && !os.IsNotExist(err) {
			return errors.New("could not read files")
		}

		for _, e := range entries {
			paths = append(paths, filepath.Join(basePath, e.Name()))
		}
	} else {
		// partial prefix
		basePath = bucketPath
		filter := fmt.Sprintf("%s/%s*", bucketPath, prefix)
		var err error
		paths, err = filepath.Glob(filter)
		if err != nil {
			return errors.New("could not filter files")
		}
	}

	hash := md5.New()
	keyCount := 0

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not read metadata")
		}

		key, err := filepath.Rel(bucketPath, path)
		if err != nil {
			return errors.New("could not resolve key")
		}

		if info.IsDir() {
			key += delimiter
			res.CommonPrefixes = append(res.CommonPrefixes, CommonPrefixes{Prefix: key})
			continue
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
