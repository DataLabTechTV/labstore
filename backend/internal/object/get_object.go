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

func GetObject(bucket, key string) (io.ReadSeekCloser, int, error) {
	objPath := filepath.Join(config.Env.StorageRoot, bucket, key)

	file, err := os.Open(objPath)
	if err != nil {
		return nil, -1, ErrorNoSuchKey()
	}

	info, err := file.Stat()
	if err != nil {
		return nil, -1, core.ErrorInternalError("Couldn't compute file size")
	}

	return file, int(info.Size()), nil
}

// GetObjectHandler: GET /:bucket/:key
func GetObjectHandler(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	key := c.Params("+")

	file, size, err := GetObject(bucket, key)
	if err != nil {
		core.HandleError(c, err)
		return err
	}
	defer file.Close()

	ext := filepath.Ext(key)
	mimeType := mime.TypeByExtension(ext)
	c.Type(mimeType)

	c.Response().Header.SetLastModified(time.Now())
	return c.SendStream(file, size)
}
