package s3

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IllumiKnowLabs/labstore/client"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

type HeadObjectResponse struct {
	ContentType   string
	ContentLength int64
	LastModified  time.Time
	StatusCode    int
}

func (c *Client) PutObject(bucket, key string, file *os.File, progress chan<- client.Progress) (int, error) {
	reqURL, err := c.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	enc := NewSigV4ChunkEncoder(c.Ctx, file, int(info.Size()), c.ChunkSize, progress)
	resp, err := c.DoSigV4Request("PUT", reqURL.String(), enc)
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

func (c *Client) HeadObject(bucket, key string) (*HeadObjectResponse, error) {
	reqURL, err := c.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return nil, err
	}

	resp, err := c.DoSigV4Request("HEAD", reqURL.String(), nil)
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

func (c *Client) GetObject(bucket, key string, writer io.Writer, progress chan<- client.Progress) (int, error) {
	reqURL, err := c.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return 0, err
	}

	resp, err := c.DoSigV4Request("GET", reqURL.String(), nil)
	if err != nil {
		return 0, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			return 0, err
		}

		buf := make([]byte, config.S3.IO.BufferSize)
		read := 0
		for {
			n, err := resp.Body.Read(buf)

			if n > 0 {
				if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
					return 0, writeErr
				}

				read += n

				if progress != nil {
					progress <- client.Progress{Current: read, Total: size}
				}
			}

			if err != nil {
				if err == io.EOF {
					break
				}
				return 0, err
			}
		}

		return resp.StatusCode, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return 0, err
	}
	s3Err.StatusCode = resp.StatusCode

	return 0, &s3Err
}

func (c *Client) DeleteObject(bucket, key string) (int, error) {
	reqURL, err := c.baseURL.Parse(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return 0, err
	}

	resp, err := c.DoSigV4Request("DELETE", reqURL.String(), nil)
	if err != nil {
		return 0, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusNoContent {
		return resp.StatusCode, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return 0, err
	}
	s3Err.StatusCode = resp.StatusCode

	return 0, &s3Err
}

func (c *Client) DeleteObjects(bucket string, keys ...string) (*types.DeleteResult, int, error) {
	if len(keys) == 1 {
		// Request DELETE /:bucket/:key
		code, err := c.DeleteObject(bucket, keys[0])
		if err != nil {
			return nil, code, err
		}
		res := &types.DeleteResult{
			Deleted: []types.DeletedObject{
				{
					DeleteMarker: false,
					Key:          keys[0],
				},
			},
		}
		return res, code, nil
	}

	reqURL, err := c.baseURL.Parse(bucket)
	if err != nil {
		return nil, 0, err
	}

	q := reqURL.Query()
	q.Set("delete", "")
	reqURL.RawQuery = q.Encode()

	var objects []types.ObjectIdentifier
	for _, key := range keys {
		object := types.ObjectIdentifier{
			BaseObject: types.BaseObject{
				Key: key,
			},
		}
		objects = append(objects, object)
	}

	req := &types.DeleteObjectsRequest{
		Object: objects,
		Quiet:  true,
	}

	xmlReader := io.NopCloser(client.XMLReader(req))
	resp, err := c.DoSigV4Request("POST", reqURL.String(), xmlReader)
	if err != nil {
		return nil, 0, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DeleteResult
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, 0, err
		}
		return &res, resp.StatusCode, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return nil, 0, err
	}
	s3Err.StatusCode = resp.StatusCode

	return nil, 0, &s3Err
}
