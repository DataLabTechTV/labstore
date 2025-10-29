package middleware

import (
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
	"github.com/gofiber/fiber/v2"
)

func WithIAM(action iam.Action, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("iamAction", action)
		return handler(c)
	}
}

func IAMMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		action, _ := c.Locals("iamAction").(string)
		if action == "" {
			return c.Next()
		}

		bucket := c.Params("bucket")
		if bucket == "" {
			return c.Next()
		}

		accessKey := c.Locals("accessKey").(string)

		if !iam.CheckPolicy(accessKey, bucket, action) {
			core.HandleError(c, core.ErrorAccessDenied())
			// !FIXME: proper S3 error handling
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
