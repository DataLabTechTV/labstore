package object

import (
	"mime"
	"path/filepath"
	"strconv"
	"time"

	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

// HeadObjectHandler: Head /:bucket/:key
func HeadObjectHandler(c *fiber.Ctx) error {
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
	c.Response().Header.Set("Content-Length", strconv.Itoa(size))

	return c.SendStatus(fiber.StatusOK)
}
