package s3

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func (c *Client) ListBuckets() ([]types.Bucket, error) {
	reqURL, err := c.baseURL.Parse("/")
	if err != nil {
		return nil, err
	}

	resp, err := c.DoSigV4Request("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var result types.ListAllMyBucketsResult
		if err := helper.ReadXML(resp.Body, &result); err != nil {
			return nil, err
		}

		return result.Buckets.Bucket, nil
	}

	var s3Err errs.S3Error
	if err := helper.ReadXML(resp.Body, &s3Err); err != nil {
		return nil, err
	}
	s3Err.StatusCode = resp.StatusCode

	return nil, &s3Err
}
