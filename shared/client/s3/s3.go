package s3

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

type S3Client struct {
	Host      string
	Port      uint16
	AccessKey string
	SecretKey string
	TLS       bool

	baseURL *url.URL
}

func NewS3Client(host string, port uint16, accessKey, secretKey string, tls bool) *S3Client {
	client := &S3Client{
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

func (client *S3Client) ListBuckets() ([]t.Bucket, error) {
	url := client.baseURL

	resp, err := client.DoSigV4Request("GET", url.String(), nil)
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

func (client *S3Client) ListObjects(bucket, key string, useV2 bool) ([]t.Object, error) {
	url, err := client.baseURL.Parse(fmt.Sprintf("%s%s", bucket, key))
	if err != nil {
		return nil, err
	}

	if useV2 {
		q := url.Query()
		q.Set("list-type", "2")
		url.RawQuery = q.Encode()
	}

	resp, err := client.DoSigV4Request("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		if useV2 {
			var result t.ListBucketResultV2
			if err := helper.ReadXML(resp.Body, &result); err != nil {
				return nil, err
			}

			return result.Contents, nil
		}

		var result t.ListBucketResult
		if err := helper.ReadXML(resp.Body, &result); err != nil {
			return nil, err
		}

		return result.Contents, nil
	}

	return []t.Object{}, nil
}
