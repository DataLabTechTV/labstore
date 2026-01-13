package s3

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/client"
)

func (client *S3Client) PutObject(bucket, key string, file *os.File, progress chan<- client.Progress) error {
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
