package iam

import (
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type CreateUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateUserResponse"`
	CreateUserResult *CreateUserResult
}

type CreateUserResult struct {
	User             *UserResult
	ResponseMetadata *ResponseMetadata
}

func (store *Store) CreateUser(name string) (*User, error) {
	if name == defaultAdminUserName {
		return nil, &errs.ErrExists{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	user := &User{
		UserID: GenerateUniqueID(IAMUserUniqueID),
		Name:   name,
		Arn:    toArn(ArnUser, defaultUserPath+name),
	}

	query := `
	INSERT INTO users (user_id, name, arn)
	VALUES (:user_id, :name, :arn)
	`

	_, err := store.sqlNamedExec(query, &user)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				slog.Warn("create user", "err", sqliteErr)
				return nil, &errs.ErrExists{Type: errs.ErrEntityTypeUser, Resource: name}
			}
		}

		slog.Error("create user", "err", err)
		return nil, err
	}

	user, err = store.GetUserByName(name)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	return user, nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.CreateUser(userName)
	if err != nil {
		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			errs.Handle(w, errs.IAMEntityAlreadyExists(errExists.Resource))
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &CreateUserResponse{
		CreateUserResult: &CreateUserResult{
			User: user.Result(),
			ResponseMetadata: &ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
