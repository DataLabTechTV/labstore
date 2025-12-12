package errs

import (
	"encoding/xml"
	"fmt"
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

func S3InternalError(message string) *S3Error {
	return &S3Error{
		Code:       "InternalError",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func S3NoSuchBucket(bucket string) *S3Error {
	return &S3Error{
		Code:       "NoSuchBucket",
		Message:    "Bucket does not exist.",
		BucketName: bucket,
		StatusCode: http.StatusNotFound,
	}
}

func S3SignatureDoesNotMatch() *S3Error {
	return &S3Error{
		Code:       "SignatureDoesNotMatch",
		Message:    "The request signature we calculate does not match the signature you provided.",
		StatusCode: http.StatusForbidden,
	}
}

func S3BucketAlreadyExists() *S3Error {
	return &S3Error{
		Code:       "BucketAlreadyExists",
		Message:    "Could not create bucket, because it already exists.",
		StatusCode: http.StatusConflict,
	}
}

func S3AccessDenied() *S3Error {
	return &S3Error{
		Code:       "AccessDenied",
		Message:    "AccessDenied",
		StatusCode: 403,
	}
}

func S3NoSuchKey(key string) *S3Error {
	return &S3Error{
		Key:        key,
		Code:       "NoSuchKey",
		Message:    "Object not found.",
		StatusCode: http.StatusNotFound,
	}
}
