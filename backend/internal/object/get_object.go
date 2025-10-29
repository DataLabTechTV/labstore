package object

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func GetObject(bucket, key string) (io.ReadSeeker, error) {
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)

	f, err := os.Open(objPath)
	if err != nil {
		return nil, ErrorNoSuchKey()
	}
	defer f.Close()

	return f, nil
}

// GetObjectHandler: GET /:bucket/:key
func GetObjectHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	key := c.Params("key")

	f, err := GetObject(bucket, key)
	if err != nil {
		core.HandleError(c, err)
		return err
	}

	ext := filepath.Ext(key)
	mimeType := mime.TypeByExtension(ext)
	c.Type(mimeType)

	c.Response().Header.SetLastModified(time.Now())
	return c.SendStream(f)
}
