package iam

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type IAMErrorResponse struct {
	XMLName   xml.Name `xml:"ErrorResponse"`
	Error     *IAMError
	RequestId string
}

type IAMError struct {
	Type       IAMErrorType
	Code       string
	Message    string
	RequestId  string `xml:"-"`
	StatusCode int    `xml:"-"`
}

type IAMErrorType string

const (
	IAMSenderType   IAMErrorType = "Sender"
	IAMReceiverType IAMErrorType = "Receiver"
)

func (e *IAMError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func ErrNotImplemented(action IAMOp) *IAMError {
	return &IAMError{
		Type:       IAMSenderType,
		Code:       "NotImplemented",
		Message:    fmt.Sprintf("The action %s is not implemented", action),
		StatusCode: http.StatusBadRequest,
	}
}

func ErrEntityAlreadyExists(entityName string) *IAMError {
	return &IAMError{
		Type:       IAMReceiverType,
		Code:       "EntityAlreadyExists",
		Message:    fmt.Sprintf("The entity %s already exists", entityName),
		StatusCode: http.StatusConflict,
	}
}

func ErrServiceFailure() *IAMError {
	return &IAMError{
		Type:       IAMReceiverType,
		Code:       "ServiceFailure",
		Message:    "The request processing has failed because of an internal error.",
		StatusCode: http.StatusInternalServerError,
	}
}
