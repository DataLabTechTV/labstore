package object

import (
	"net/http"

	"github.com/DataLabTechTV/labstore/backend/internal/core"
)

func ErrorNoSuchKey() *core.S3Error {
	return &core.S3Error{
		Code:       "NoSuchKey",
		Message:    "Object not found",
		StatusCode: http.StatusNotFound,
	}
}
