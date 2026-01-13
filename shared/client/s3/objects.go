package s3

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/client"
)

type HeadObjectResponse struct {
	ContentType   string
	ContentLength int64
	LastModified  time.Time
	StatusCode    int
}

func (client *S3Client) PutObject(bucket, key string, file *os.File, progress chan<- client.Progress) (int, error) {
	reqURL, err := client.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	enc := NewSigV4ChunkEncoder(file, int(info.Size()), client.ChunkSize, progress)
	resp, err := client.DoSigV4Request("PUT", reqURL.String(), enc)
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

func (client *S3Client) HeadObject(bucket, key string) (*HeadObjectResponse, error) {
	reqURL, err := client.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return nil, err
	}

	resp, err := client.DoSigV4Request("HEAD", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}

	lastModified, err := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		out := &HeadObjectResponse{
			ContentType:   resp.Header.Get("Content-Type"),
			ContentLength: contentLength,
			LastModified:  lastModified,
			StatusCode:    resp.StatusCode,
		}
		return out, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return nil, err
	}
	s3Err.StatusCode = resp.StatusCode

	return nil, &s3Err
}
