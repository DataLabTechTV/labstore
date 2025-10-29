package bucket

import (
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func DeleteBucket(bucket string) error {
	path := filepath.Join(config.Env.StorageRoot, bucket)

	err := os.Remove(path)
	if err != nil {
		return core.ErrorNoSuchBucket()
	}

	return nil
}

// DeleteBucketHandler: DELETE /:bucket
func DeleteBucketHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	if err := DeleteBucket(bucket); err != nil {
		core.HandleError(c, err)
		return err
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
