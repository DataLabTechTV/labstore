package errs

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
)

const (
	ErrEntityTypeBucket      = "bucket"
	ErrEntityTypeObject      = "object"
	ErrEntityTypeUser        = "user"
	ErrEntityTypeAccessKey   = "access_key"
	ErrEntityTypeGroup       = "group"
	ErrEntityTypeUserGroup   = "user_group"
	ErrEntityTypePolicy      = "policy"
	ErrEntityTypeUserPolicy  = "user_policy"
	ErrEntityTypeGroupPolicy = "group_policy"
)

type ErrEntityType string

type ErrExists struct {
	Type     ErrEntityType
	Resource string
}

type ErrNotFound struct {
	Type     ErrEntityType
	Resource string
}

type ErrForbidden struct {
	Type     ErrEntityType
	Resource string
}

type ErrUserPolicyNotAttached struct {
	UserID   string
	PolicyID string
}

type ErrGroupPolicyNotAttached struct {
	GroupID  string
	PolicyID string
}

type ErrUserNotInGroup struct {
	UserID  string
	GroupID string
}

type ErrPathTraversal struct {
	BaseDir string
	Path    string
}

type ErrNotDirectory struct {
	Path string
}

func (e *ErrExists) Error() string {
	return fmt.Sprintf("%s exists: %s", e.Type, e.Resource)
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Type, e.Resource)
}

func (e *ErrForbidden) Error() string {
	return fmt.Sprintf("Operation forbidden for %s: %s", e.Type, e.Resource)
}

func (e *ErrUserPolicyNotAttached) Error() string {
	return fmt.Sprintf("Policy %s not attached to user %s", e.PolicyID, e.UserID)
}

func (e *ErrUserPolicyNotAttached) Description() string {
	return fmt.Sprintf("{UserID=%s, PolicyID=%s}", e.UserID, e.PolicyID)
}

func (e *ErrGroupPolicyNotAttached) Error() string {
	return fmt.Sprintf("Policy %s not attached to group %s", e.PolicyID, e.GroupID)
}

func (e *ErrGroupPolicyNotAttached) Description() string {
	return fmt.Sprintf("{GroupID=%s, PolicyID=%s}", e.GroupID, e.PolicyID)
}

func (e *ErrUserNotInGroup) Error() string {
	return fmt.Sprintf("User %s does not belong to group %s", e.UserID, e.GroupID)
}

func (e *ErrUserNotInGroup) Description() string {
	return fmt.Sprintf("{UserID=%s, GroupID=%s}", e.UserID, e.GroupID)
}

func (e *ErrPathTraversal) Error() string {
	return fmt.Sprintf("Invalid path traversal {BaseDir=%s, Path=%s}", e.BaseDir, e.Path)
}

func (e *ErrNotDirectory) Error() string {
	return fmt.Sprintf("Not a directory: %s", e.Path)
}

func Handle(w http.ResponseWriter, err error) {
	var s3Err *S3Error
	var iamErr *IAMError

	if errors.As(err, &s3Err) {
		slog.Error("s3 error", "err", s3Err)
		core.WriteXML(w, s3Err.StatusCode, s3Err)
	} else if errors.As(err, &iamErr) {
		slog.Error("iam error", "err", iamErr)
		core.WriteXML(w, iamErr.StatusCode, IAMErrorResponse{
			Error:     iamErr,
			RequestId: iamErr.RequestId,
		})
	} else {
		slog.Error("internal server error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
