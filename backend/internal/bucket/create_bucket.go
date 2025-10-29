package bucket

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func ErrorBucketAlreadyExists() *core.S3Error {
	return &core.S3Error{
		Code:       "BucketAlreadyExists",
		Message:    "Could not create bucket, because it already exists",
		StatusCode: fiber.StatusConflict,
	}
}

func CreateBucket(bucket string) error {
	path := filepath.Join(config.Env.StorageRoot, bucket)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ErrorBucketAlreadyExists()
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("could not create bucket: %w", err)
	}

	return nil
}

// CreateBucket: PUT /:bucket
func PutBucketHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")

	if err := CreateBucket(bucket); err != nil {
		core.HandleError(c, err)
		return err
	}

	c.Status(fiber.StatusOK)
	return nil
}
