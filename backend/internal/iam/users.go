package iam

import (
	"errors"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type User struct {
	Name        string `db:"name"`
	AccessKeyID string
	SecretKey   string

	GroupIDs  []string
	PolicyIDs []string
}

func GetUser(accessKey string) (*User, bool) {
	user, ok := store.Users[accessKey]
	return user, ok
}

func CreateUser(name string) *errs.IAMError {
	if name == config.Admin.Auth.AccessKey {
		return errs.IAMEntityAlreadyExists(name)
	}

	user := User{Name: name}

	_, err := store.writeDB.NamedExec(`INSERT INTO users (name) VALUES (:name)`, &user)
	if err != nil {
		return errs.IAMServiceFailure()
	}

	return nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errors.New("missing query parameter: UserName"))
		return
	}

	if err := CreateUser(userName); err != nil {
		errs.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (store *Store) setupAdmin() {
	store.Users[config.Admin.Auth.AccessKey] = &User{
		Name:        "Administrator",
		AccessKeyID: config.Admin.Auth.AccessKey,
		SecretKey:   config.Admin.Auth.SecretKey,
		PolicyIDs:   []string{adminPolicy},
	}

	store.Policies[adminPolicy] = &Policy{
		ID: adminPolicy,
		Document: &PolicyDocument{
			Version: latestPolicyDocumentVersion,
			Statement: []Statement{
				{
					Effect:    allow,
					Actions:   []Action{Action(Any)},
					Resources: []string{Any},
				},
			},
		},
	}
}
