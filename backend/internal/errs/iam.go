package errs

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
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

func IAMNotImplemented(action string) *IAMError {
	return &IAMError{
		Type:       IAMSenderType,
		Code:       "NotImplemented",
		Message:    fmt.Sprintf("The action %s is not implemented.", action),
		StatusCode: http.StatusBadRequest,
	}
}

func IAMServiceFailure() *IAMError {
	return &IAMError{
		Type:       IAMReceiverType,
		Code:       "ServiceFailure",
		Message:    "The request processing has failed because of an internal error.",
		StatusCode: http.StatusInternalServerError,
	}
}

func IAMEntityAlreadyExists(entityName string) *IAMError {
	return &IAMError{
		Type:       IAMReceiverType,
		Code:       "EntityAlreadyExists",
		Message:    fmt.Sprintf("The entity %s already exists.", entityName),
		StatusCode: http.StatusConflict,
	}
}

func IAMNoSuchEntity(entityType, entityName string) *IAMError {
	return &IAMError{
		Type: IAMReceiverType,
		Code: "NoSuchEntity",
		Message: fmt.Sprintf(
			"The %s %s does not exist.",
			strings.ReplaceAll(entityType, "_", " "),
			entityName,
		),
		StatusCode: http.StatusNotFound,
	}
}

func IAMDeleteConflict(entityName string) *IAMError {
	return &IAMError{
		Type:       IAMReceiverType,
		Code:       "DeleteConflict",
		Message:    fmt.Sprintf("Cannot delete the entity %s.", entityName),
		StatusCode: http.StatusConflict,
	}
}
