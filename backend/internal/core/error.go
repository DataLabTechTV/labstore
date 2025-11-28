package core

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type S3Error struct {
	XMLName    xml.Name `xml:"Error"`
	Code       string
	Message    string
	Resource   string
	BucketName string
	Key        string
	VersionId  string
	RequestId  string
	HostId     string
	StatusCode int `xml:"-"`
}

func (e *S3Error) Error() string {
	if e.Key == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}

	return fmt.Sprintf("%s: %s: %s", e.Code, e.Message, e.Key)
}

func (e *S3Error) WithRequestID(requestID string) *S3Error {
	e.RequestId = requestID
	return e
}

func (e *S3Error) WithHostID(hostID string) *S3Error {
	e.HostId = hostID
	return e
}

func (e *S3Error) WithResource(resource string) *S3Error {
	e.Resource = resource
	return e
}

func HandleError(w http.ResponseWriter, err error) {
	var s3Err *S3Error

	if errors.As(err, &s3Err) {
		slog.Error("s3 error", "err", s3Err)
		WriteXML(w, s3Err.StatusCode, s3Err)
	} else {
		slog.Error("internal server error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrInternalError(message string) *S3Error {
	return &S3Error{
		Code:       "InternalError",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func ErrNoSuchBucket(bucket string) *S3Error {
	return &S3Error{
		Code:       "NoSuchBucket",
		Message:    "Bucket does not exist",
		BucketName: bucket,
		StatusCode: http.StatusNotFound,
	}
}
