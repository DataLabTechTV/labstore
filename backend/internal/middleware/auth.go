package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/DataLabTechTV/labstore/backend/internal/auth"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
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
			// logger.Log.Error(err.Error()) // if we need to save error to logs
			core.HandleError(c, core.ErrorSignatureDoesNotMatch())
			return c.Next()
		}

		c.Locals("accessKey", accessKey)

		return c.Next()
	}
}
