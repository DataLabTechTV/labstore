package bucket

import (
	"github.com/gofiber/fiber/v2"
)

func HeadBucketHandler(c *fiber.Ctx) error {
	// *NOTE: This will likely share code with GET due to using the same headers.
	// TODO: organize shared code somewhere

	c.Status(fiber.StatusOK)
	return nil
}
