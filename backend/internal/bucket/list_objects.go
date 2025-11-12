package bucket

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
	"github.com/DataLabTechTV/labstore/backend/internal/object"
)

const DefaultMaxKeys = 1000

type ListObjectsRequest struct {
	Bucket  string
	Prefix  string
	MaxKeys int
	// TODO: ...
}

type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	IsTruncated bool
	Contents    []object.Object
}

func ListObjects(r *ListObjectsRequest) (*ListBucketResult, error) {
	slog.Debug("Processing ListObjects")

	res := &ListBucketResult{
		Name:        r.Bucket,
		MaxKeys:     r.MaxKeys,
		Contents:    []object.Object{},
		IsTruncated: false,
	}

	hash := md5.New()
	keyCount := 0
	basePath := filepath.Join(config.Env.StorageRoot, r.Bucket, r.Prefix)

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return err
		}

		key, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		lastModified := core.Timestamp(info.ModTime())

		if _, err := io.Copy(hash, file); err != nil {
			return err
		}
		eTag := hex.EncodeToString(hash.Sum(nil))

		size := info.Size()

		obj := object.Object{
			Key:          key,
			LastModified: lastModified,
			ETag:         eTag,
			Size:         size,
			// TODO add missing Owner when there is proper IAM
		}

		res.Contents = append(res.Contents, obj)

		if keyCount++; keyCount >= res.MaxKeys {
			// res.IsTruncated = true
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, errors.New("failed to list objects")
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
			Bucket:  bucket,
			Prefix:  prefix,
			MaxKeys: maxKeys,
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
