package s3

import (
	"fmt"
	"net/http"
	"strings"

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

	uri string
}

func NewS3Client(host string, port uint16, accessKey, secretKey string, tls bool) *S3Client {
	client := &S3Client{
		Host:      host,
		Port:      port,
		AccessKey: accessKey,
		SecretKey: secretKey,
		TLS:       tls,
	}

	var uri strings.Builder
	if tls {
		uri.WriteString("http://")
	} else {
		uri.WriteString("http://")
	}

	uri.WriteString(client.Host)
	uri.WriteRune(':')
	uri.WriteString(fmt.Sprintf("%d", client.Port))

	client.uri = uri.String()

	return client
}

func (client *S3Client) ListBuckets() ([]string, error) {
	resp, err := http.Get(client.uri)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var result t.ListAllMyBucketsResult
		if err := helper.ReadXML(resp.Body, &result); err != nil {
			return nil, err
		}

		var list []string
		for _, bucket := range result.Buckets.Bucket {
			list = append(list, bucket.Name)
		}

		return list, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return nil, err
	}

	return nil, &s3Err
}

func (client *S3Client) ListObjects(bucket, key string) ([]string, error) {
	// TODO: implement
	return []string{}, nil
}
