package bucket

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

const DefaultDelimiter = "/"

type BaseListObjectsRequest struct {
	Bucket    string
	Prefix    string
	Delimiter string
	MaxKeys   int
	afterKey  string
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

// ListObjectsHandler: GET /:bucket
func ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	var res any
	var err error
	var delimiter string
	var maxKeys int

	bucket := r.PathValue("bucket")
	requestID := core.NewRequestID()

	q := r.URL.Query()

	prefix := q.Get("prefix")

	if d := q.Get("delimiter"); d == "" {
		delimiter = DefaultDelimiter
	} else {
		delimiter = d
	}

	if mk := q.Get("maxKeys"); mk == "" {
		maxKeys = config.S3.Paging.MaxKeys
	} else {
		if maxKeys, err = strconv.Atoi(mk); err != nil {
			slog.Warn("invalid max-keys value, using default", "input", mk, "default", config.S3.Paging.MaxKeys)
			maxKeys = config.S3.Paging.MaxKeys
		}
	}

	if maxKeys > config.S3.Paging.MaxKeys {
		slog.Warn("max-keys capped", "input", maxKeys, "cap", config.S3.Paging.MaxKeys)
		maxKeys = config.S3.Paging.MaxKeys
	}

	rBase := BaseListObjectsRequest{
		Bucket:    bucket,
		Prefix:    prefix,
		Delimiter: delimiter,
		MaxKeys:   maxKeys,
	}

	if !core.BucketExists(rBase.Bucket) {
		errs.Handle(w, errs.S3NoSuchBucket(rBase.Bucket))
		return
	}

	if rBase.Delimiter != "/" {
		errs.Handle(w, errors.New("only '/' delimiters are supported by LabStore"))
		return
	}

	if q.Get("list-type") == "2" {
		continuationToken := q.Get("continuation-token")
		startAfter := q.Get("start-after")
		fetchOwner := q.Get("fetch-owner") == "true"

		var token []byte
		token, err = base64.StdEncoding.DecodeString(continuationToken)
		if err != nil {
			errs.Handle(w, errs.S3InternalError())
			return
		}
		rBase.afterKey = string(token)

		r := &ListObjectsRequestV2{
			BaseListObjectsRequest: rBase,
			ContinuationToken:      continuationToken,
			StartAfter:             startAfter,
			FetchOwner:             fetchOwner,
		}

		res, err = ListObjectsV2(r)
	} else {
		marker := q.Get("marker")
		rBase.afterKey = marker

		r := &ListObjectsRequest{
			BaseListObjectsRequest: rBase,
			Marker:                 marker,
		}

		res, err = ListObjects(r)
	}

	if err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			slog.Error("list objects", "err", errForbidden)
			errs.Handle(w, errs.S3AccessDenied())
		}

		slog.Error("list objects", "err", err)
		errs.Handle(w, errs.S3InternalError())
		return
	}

	w.Header().Set("Server", "LabStore")
	w.Header().Set("X-Amz-Request-Id", requestID)

	helper.WriteXMLResponse(w, http.StatusOK, res)
}

func ListObjects(r *ListObjectsRequest) (*t.ListBucketResult, error) {
	slog.Debug("list objects", "request", r)

	res := &t.ListBucketResult{
		BaseListBucketResult: t.BaseListBucketResult{
			Name:        r.Bucket,
			MaxKeys:     r.MaxKeys,
			IsTruncated: false,
		},
	}

	err := list(&res.BaseListBucketResult, &r.BaseListObjectsRequest)
	if err != nil {
		return nil, err
	}

	res.MaxKeys = r.MaxKeys
	res.Marker = r.Marker
	res.NextMarker = res.UntilKey

	return res, nil
}

func ListObjectsV2(r *ListObjectsRequestV2) (*t.ListBucketResultV2, error) {
	slog.Debug("list objects v2", "request", r)

	res := &t.ListBucketResultV2{
		BaseListBucketResult: t.BaseListBucketResult{
			Name:        r.Bucket,
			MaxKeys:     r.MaxKeys,
			IsTruncated: false,
		},
	}

	err := list(&res.BaseListBucketResult, &r.BaseListObjectsRequest)
	if err != nil {
		return nil, err
	}

	res.MaxKeys = r.MaxKeys
	res.StartAfter = r.StartAfter
	res.ContinuationToken = r.ContinuationToken
	res.NextContinuationToken = base64.StdEncoding.EncodeToString([]byte(res.UntilKey))
	res.KeyCount = len(res.Contents)

	return res, nil
}

// Lists objects as Contents, and directories as CommonPrefixes, for a given fs path
func list(res *t.BaseListBucketResult, req *BaseListObjectsRequest) error {
	bucketPath := core.BucketSystemPath(req.Bucket)

	if !security.IsSubdir(config.Storage.ObjectsPath, bucketPath) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeBucket, Resource: req.Bucket}
	}

	filterPath := fmt.Sprintf("%s%c%s*", bucketPath, os.PathSeparator, req.Prefix)

	if !security.IsSubdir(config.Storage.ObjectsPath, filterPath) {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeObject, Resource: core.ObjectPath(req.Bucket, req.Prefix)}
	}

	paths, err := filepath.Glob(filterPath)
	if err != nil {
		return errors.New("could not filter files")
	}

	slices.Sort(paths)

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

		if req.afterKey > key {
			continue
		}

		if info.IsDir() {
			key += req.Delimiter
			res.CommonPrefixes = append(res.CommonPrefixes, t.CommonPrefix{Prefix: key})
			continue
		}

		lastModified := t.Timestamp(info.ModTime())

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not read file: %s", key)
		}
		defer helper.CloseWithErr(file, &err)

		if _, err := io.Copy(hash, file); err != nil {
			return fmt.Errorf("could not compute hash: %s", key)
		}
		eTag := hex.EncodeToString(hash.Sum(nil))

		size := info.Size()

		obj := t.Object{
			BaseObject: t.BaseObject{
				Key:          key,
				LastModified: lastModified,
				ETag:         eTag,
				Size:         size,
			},
			// TODO: ...
		}

		res.Contents = append(res.Contents, obj)

		if keyCount++; keyCount > res.MaxKeys {
			res.UntilKey = key
			res.IsTruncated = true
			return nil
		}
	}

	return nil
}
