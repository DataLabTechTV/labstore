package bucket

import (
	"encoding/xml"

	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	IsTruncated bool
	Marker      string
	Name        string
	Prefix      string
	MaxKeys     int
}

func ListBucket(bucket string) (*ListBucketResult, error) {
	// TODO: implement
	return &ListBucketResult{}, nil
}

// ListObjectsHandler: GET /:bucket
func ListObjectsHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	requestID := middleware.NewRequestID()

	logger.Log.Info("Received probe: ", c.Path())

	res, err := ListBucket(bucket)
	if err != nil {
		core.HandleError(c, err)
		return err
	}

	c.Set("Server", "LabStore")
	c.Set("X-Amz-Request-Id", requestID)

	return c.XML(res)
}
