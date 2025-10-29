package object

import (
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	"github.com/gofiber/fiber/v2"
)

func HeadObject(bucket, key string) error {
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)

	if !helper.FileExists(objPath) {
		return ErrorNoSuchKey()
	}

	return nil
}

// HeadObjectHandler: Head /:bucket/:key
func HeadObjectHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	key := c.Params("key")

	if err := HeadObject(bucket, key); err != nil {
		core.HandleError(c, err)
		return err
	}

	c.Status(fiber.StatusOK)
	return nil
}
