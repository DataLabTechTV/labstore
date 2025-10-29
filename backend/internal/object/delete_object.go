package object

import (
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func DeleteObject(bucket, key string) error {
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)

	err := os.Remove(objPath)
	if err != nil {
		return ErrorNoSuchKey()
	}

	return nil
}

// DeleteObjectHandler: DELETE /:bucket/:key
func DeleteObjectHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	key := c.Params("+")

	if err := DeleteObject(bucket, key); err != nil {
		core.HandleError(c, err)
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
