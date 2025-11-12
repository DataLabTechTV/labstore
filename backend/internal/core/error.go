package core

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
)

func ErrorAccessDenied() *S3Error {
	return &S3Error{
		Code:       "AccessDenied",
		Message:    "AccessDenied",
		StatusCode: 403,
	}
}

func ErrorNotImplemented() *S3Error {
	return &S3Error{
		Code:       "NotImplemented",
		Message:    "Operation not implemented",
		StatusCode: fiber.StatusNotImplemented,
	}
}

func ErrorInternalError(message string) *S3Error {
	return &S3Error{
		Code:       "InternalError",
		Message:    message,
		StatusCode: fiber.StatusInternalServerError,
	}
}

func ErrorSignatureDoesNotMatch() *S3Error {
	return &S3Error{
		Code:		"SignatureDoesNotMatch",
		Message:	"The request signature does not match the signature you provide",
		StatusCode:	403,
	}
}

func ErrorNoSuchBucket() *S3Error {
	return &S3Error{
		Code:       "NoSuckBucket",
		Message:    "Bucket does not exist",
		StatusCode: fiber.StatusNotFound,
	}
}

type S3Error struct {
	XMLName    xml.Name `xml:"Error"`
	Code       string
	Message    string
	RequestId  string
	HostId     string
	StatusCode int `xml:"-"`
}

func (e *S3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *S3Error) WithRequestID(requestID string) *S3Error {
	e.RequestId = requestID
	return e
}

func (e *S3Error) WithHostID(hostID string) *S3Error {
	e.HostId = hostID
	return e
}

func HandleError(c *fiber.Ctx, err error) {
	logger.Log.Error(err.Error())

	var s3Error *S3Error

	if errors.As(err, &s3Error) {
		c.Status(s3Error.StatusCode).XML(s3Error)
	} else {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
}
