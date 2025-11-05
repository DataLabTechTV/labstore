package core

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

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
		StatusCode: http.StatusNotImplemented,
	}
}

func ErrorInternalError(message string) *S3Error {
	return &S3Error{
		Code:       "InternalError",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func ErrorNoSuchBucket() *S3Error {
	return &S3Error{
		Code:       "NoSuckBucket",
		Message:    "Bucket does not exist",
		StatusCode: http.StatusNotFound,
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

func HandleError(w http.ResponseWriter, err error) {
	logger.Log.Errorf("Server error: %s", err.Error())

	var s3Error *S3Error

	if errors.As(err, &s3Error) {
		WriteXML(w, s3Error.StatusCode, s3Error)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
