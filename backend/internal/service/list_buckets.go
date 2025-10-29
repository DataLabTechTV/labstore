package service

import (
	"encoding/xml"
	"os"
	"time"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

// !FIXME: move types to a proper location

type Bucket struct {
	Name         string
	CreationDate string
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   struct {
		ID          string
		DisplayName string
	}
	Buckets struct {
		Bucket []Bucket
	}
}

func ListBuckets(accessKey string) (*ListAllMyBucketsResult, error) {
	entries, err := os.ReadDir(config.Env.StorageRoot)
	if err != nil {
		return nil, core.ErrorInternalError("Failed to list buckets")
	}

	res := ListAllMyBucketsResult{}
	res.Owner.ID = accessKey
	res.Owner.DisplayName = accessKey

	for _, e := range entries {
		if e.IsDir() {
			b := Bucket{Name: e.Name(), CreationDate: time.Now().Format(time.RFC3339)}
			res.Buckets.Bucket = append(res.Buckets.Bucket, b)
		}
	}

	return &res, nil
}

// ListBuckets: GET /
func ListBucketsHandler(c *fiber.Ctx) error {
	accessKey := c.Params("accessKey")

	res, err := ListBuckets(accessKey)
	if err != nil {
		core.HandleError(c, err)
		return err
	}

	return c.XML(res)
}
