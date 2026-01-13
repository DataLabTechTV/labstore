package s3

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func (client *S3Client) ListBuckets() ([]types.Bucket, error) {
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

	return nil, &s3Err
}
