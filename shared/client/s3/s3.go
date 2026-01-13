package s3

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/IllumiKnowLabs/labstore/client/types"
)

type S3Client struct {
	Ctx       context.Context
	Host      string
	Port      uint16
	AccessKey string
	SecretKey string
	TLS       bool
	ChunkSize int

	baseURL *url.URL
}

type ListResult struct {
	Object       *t.Object
	CommonPrefix *t.CommonPrefix
	Err          error
}

func (lr *ListResult) IsObject() bool {
	return lr.Object != nil && lr.CommonPrefix == nil
}

func (lr *ListResult) IsCommonPrefix() bool {
	return lr.CommonPrefix != nil && lr.Object == nil
}

func NewS3Client(
	ctx context.Context,
	host string,
	port uint16,
	accessKey string,
	secretKey string,
	tls bool,
	chunkSize int,
) *S3Client {
	client := &S3Client{
		Ctx:       ctx,
		Host:      host,
		Port:      port,
		AccessKey: accessKey,
		SecretKey: secretKey,
		TLS:       tls,
		ChunkSize: chunkSize,
	}

	var scheme string
	if tls {
		scheme = "http"
	} else {
		scheme = "http"
	}

	client.baseURL = &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", client.Host, client.Port),
	}

	return client
}

func (client *S3Client) IsDone() bool {
	select {
	case <-client.Ctx.Done():
		return true
	default:
		return false
	}
}

func (client *S3Client) ListBuckets() ([]t.Bucket, error) {
	reqURL, err := client.baseURL.Parse("/")
	if err != nil {
		return nil, err
	}

	resp, err := client.DoSigV4Request("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var result t.ListAllMyBucketsResult
		if err := helper.ReadXML(resp.Body, &result); err != nil {
			return nil, err
		}

		return result.Buckets.Bucket, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return nil, err
	}

	return nil, &s3Err
}

func (client *S3Client) ListObjects(bucket, key string, useV2 bool) <-chan ListResult {
	out := make(chan ListResult)

	go func() {
		reqURL, err := client.baseURL.Parse(bucket)
		if err != nil {
			out <- ListResult{Err: err}
			close(out)
			return
		}

		q := reqURL.Query()

		q.Set("prefix", key)

		if useV2 {
			q.Set("list-type", "2")
		}

		reqURL.RawQuery = q.Encode()

		for {
			if client.IsDone() {
				close(out)
				return
			}

			slog.Debug("list objects", "reqURL", reqURL)
			resp, err := client.DoSigV4Request("GET", reqURL.String(), nil)
			if err != nil {
				out <- ListResult{Err: err}
				close(out)
				return
			}
			defer helper.CloseWithErr(resp.Body, &err)

			if resp.StatusCode != http.StatusOK {
				var s3Error errs.S3Error
				if err := helper.ReadXML(resp.Body, &s3Error); err != nil {
					out <- ListResult{Err: err}
					return
				}
				out <- ListResult{Err: &s3Error}
				close(out)
				return
			}

			if useV2 {
				var r t.ListBucketResultV2
				if err := helper.ReadXML(resp.Body, &r); err != nil {
					out <- ListResult{Err: err}
					close(out)
					return
				}

				for _, commonPrefix := range r.CommonPrefixes {
					out <- ListResult{CommonPrefix: &commonPrefix}
				}

				for _, object := range r.Contents {
					out <- ListResult{Object: &object}
				}

				if r.NextContinuationToken == "" {
					close(out)
					return
				}

				q := reqURL.Query()
				q.Set("continuation-token", r.ContinuationToken)
				reqURL.RawQuery = q.Encode()
			} else {
				var r t.ListBucketResult
				if err := helper.ReadXML(resp.Body, &r); err != nil {
					out <- ListResult{Err: err}
					close(out)
					return
				}

				for _, commonPrefix := range r.CommonPrefixes {
					out <- ListResult{CommonPrefix: &commonPrefix}
				}

				for _, object := range r.Contents {
					out <- ListResult{Object: &object}
				}

				if !r.IsTruncated {
					close(out)
					return
				}

				q := reqURL.Query()
				q.Set("marker", r.NextMarker)
				reqURL.RawQuery = q.Encode()
			}
		}
	}()

	return out
}

func (client *S3Client) PutObject(bucket, key string, file *os.File, progress chan<- types.Progress) error {
	reqURL, err := client.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	enc := NewSigV4ChunkEncoder(file, int(info.Size()), client.ChunkSize, progress)
	resp, err := client.DoSigV4Request("PUT", reqURL.String(), enc)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return err
	}

	return &s3Err

}
