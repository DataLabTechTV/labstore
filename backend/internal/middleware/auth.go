package middleware

import (
	"github.com/DataLabTechTV/labstore/backend/internal/auth"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

var ErrorInvalidAccessKey = &core.S3Error{
	Code:       "InvalidAccessKeyId",
	Message:    "Signature or access key invalid",
	StatusCode: fiber.StatusForbidden,
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessKey, err := auth.VerifyAWSSigV4(c)
		if err != nil {
			core.HandleError(c, err)
			return c.Next()
		}

		c.Locals("accessKey", accessKey)

		return c.Next()
	}
}
