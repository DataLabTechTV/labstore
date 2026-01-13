package s3

import (
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

type ListResult struct {
	Object       *types.Object
	CommonPrefix *types.CommonPrefix
	Err          error
}

func (lr *ListResult) IsObject() bool {
	return lr.Object != nil && lr.CommonPrefix == nil
}

func (lr *ListResult) IsCommonPrefix() bool {
	return lr.CommonPrefix != nil && lr.Object == nil
}

func (client *S3Client) CreateBucket(bucket string) (int, error) {
	reqURL, err := client.baseURL.Parse(bucket)
	if err != nil {
		return 0, err
	}

	resp, err := client.DoSigV4Request("PUT", reqURL.String(), nil)
	if err != nil {
		return 0, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		return resp.StatusCode, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return 0, err
	}
	s3Err.StatusCode = resp.StatusCode

	return 0, &s3Err
}

func (client *S3Client) HeadBucket(bucket string) (int, error) {
	reqURL, err := client.baseURL.Parse(bucket)
	if err != nil {
		return 0, err
	}

	resp, err := client.DoSigV4Request("HEAD", reqURL.String(), nil)
	if err != nil {
		return 0, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	return resp.StatusCode, nil
}

func (client *S3Client) ListObjects(bucket, key string, useV2 bool) <-chan ListResult {
	out := make(chan ListResult)

	go func() {
		defer close(out)

		reqURL, err := client.baseURL.Parse(bucket)
		if err != nil {
			out <- ListResult{Err: err}
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
				return
			}

			slog.Debug("list objects", "reqURL", reqURL)
			resp, err := client.DoSigV4Request("GET", reqURL.String(), nil)
			if err != nil {
				out <- ListResult{Err: err}
				return
			}
			defer helper.CloseWithErr(resp.Body, &err)

			if resp.StatusCode != http.StatusOK {
				var s3Error errs.S3Error
				if err := helper.ReadXML(resp.Body, &s3Error); err != nil {
					out <- ListResult{Err: err}
					return
				}
				s3Error.StatusCode = resp.StatusCode
				out <- ListResult{Err: &s3Error}
				return
			}

			if useV2 {
				var r types.ListBucketResultV2
				if err := helper.ReadXML(resp.Body, &r); err != nil {
					out <- ListResult{Err: err}
					return
				}

				for _, commonPrefix := range r.CommonPrefixes {
					out <- ListResult{CommonPrefix: &commonPrefix}
				}

				for _, object := range r.Contents {
					out <- ListResult{Object: &object}
				}

				if r.NextContinuationToken == "" {
					return
				}

				q := reqURL.Query()
				q.Set("continuation-token", r.ContinuationToken)
				reqURL.RawQuery = q.Encode()
			} else {
				var r types.ListBucketResult
				if err := helper.ReadXML(resp.Body, &r); err != nil {
					out <- ListResult{Err: err}
					return
				}

				for _, commonPrefix := range r.CommonPrefixes {
					out <- ListResult{CommonPrefix: &commonPrefix}
				}

				for _, object := range r.Contents {
					out <- ListResult{Object: &object}
				}

				if !r.IsTruncated {
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
