package object

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func PutObject(bucket string, key string, data []byte) error {
	bucketPath := filepath.Join(config.Env.StorageRoot, bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return core.ErrorNoSuchBucket()
	}

	objPath := filepath.Join(bucketPath, key)
	objDir := filepath.Dir(objPath)
	os.MkdirAll(objDir, 0755)

	f, err := os.Create(objPath)
	if err != nil {
		return core.ErrorInternalError("Failed to create object")
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(data))
	if err != nil {
		return core.ErrorInternalError("Failed to write object")
	}

	return nil
}

// PutObjectHandler: PUT /:bucket/:key
func PutObjectHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	key := c.Params("+")
	data := c.Body()

	if err := PutObject(bucket, key, data); err != nil {
		core.HandleError(c, err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
