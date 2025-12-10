package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
)

type User struct {
	Name        string
	AccessKeyID string
	SecretKey   string

	GroupIDs  []string
	PolicyIDs []string
}

func GetUser(accessKey string) (*User, bool) {
	user, ok := store.Users[accessKey]
	return user, ok
}

func CreateUser(username string) *IAMError {
	if username == config.Admin.Auth.AccessKey {
		return ErrEntityAlreadyExists(username)
	}

	return nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func setupAdmin() {
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
