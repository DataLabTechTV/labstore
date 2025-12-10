package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
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

func CreateUser(name string) *IAMError {
	if name == config.Admin.Auth.AccessKey {
		return ErrEntityAlreadyExists(name)
	}

	user := User{Name: name}

	_, err := store.writeDB.NamedExec(`INSERT INTO users (name) VALUES (:name)`, &user)
	if err != nil {
		return ErrServiceFailure()
	}

	return nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: call CreateUser
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
