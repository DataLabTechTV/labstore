package s3

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

type S3Client struct {
	Ctx       context.Context
	Host      string
	Port      uint16
	AccessKey string
	SecretKey string
	TLS       bool

	baseURL *url.URL
}

type Result[T any] struct {
	Value T
	Err   error
}

func NewS3Client(ctx context.Context, host string, port uint16, accessKey, secretKey string, tls bool) *S3Client {
	client := &S3Client{
		Ctx:       ctx,
		Host:      host,
		Port:      port,
		AccessKey: accessKey,
		SecretKey: secretKey,
		TLS:       tls,
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

func (client *S3Client) ListObjects(bucket, key string, useV2 bool) <-chan Result[t.Object] {
	out := make(chan Result[t.Object])

	go func() {
		reqURL, err := client.baseURL.Parse(bucket)
		if err != nil {
			out <- Result[t.Object]{Err: err}
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
				out <- Result[t.Object]{Err: err}
				close(out)
				return
			}
			defer helper.CloseWithErr(resp.Body, &err)

			if resp.StatusCode != http.StatusOK {
				var s3Error errs.S3Error
				if err := helper.ReadXML(resp.Body, &s3Error); err != nil {
					out <- Result[t.Object]{Err: err}
					return
				}
				out <- Result[t.Object]{Err: &s3Error}
				close(out)
				return
			}

			if useV2 {
				var r t.ListBucketResultV2
				if err := helper.ReadXML(resp.Body, &r); err != nil {
					out <- Result[t.Object]{Err: err}
					close(out)
					return
				}

				for _, object := range r.Contents {
					out <- Result[t.Object]{Value: object}
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
					out <- Result[t.Object]{Err: err}
					close(out)
					return
				}

				for _, object := range r.Contents {
					out <- Result[t.Object]{Value: object}
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
