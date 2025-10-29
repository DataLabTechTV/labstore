package object

import (
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/gofiber/fiber/v2"
)

func ErrorNoSuchKey() *core.S3Error {
	return &core.S3Error{
		Code:       "NoSuchKey",
		Message:    "Object not found",
		StatusCode: fiber.StatusNotFound,
	}
}
