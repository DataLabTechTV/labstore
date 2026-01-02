package iam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (store *Store) CreateUser(ctx context.Context, name string) (*User, error) {
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

	_, err := store.sqlNamedExecContext(ctx, query, &user)
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

	user, err = store.GetUserByName(ctx, name)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	return user, nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	user, err := store.CreateUser(ctx, userName)
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

	response := &t.CreateUserResponse{
		CreateUserResult: &t.CreateUserResult{
			User: user.Result(),
			ResponseMetadata: &t.ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
