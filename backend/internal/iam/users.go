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

func CreateUser() {

}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func createAdmin() {
	store.Users = map[string]*User{
		config.Admin.Auth.AccessKey: {
			Name:        "Administrator",
			AccessKeyID: "admin",
			SecretKey:   config.Admin.Auth.SecretKey,
			PolicyIDs:   []string{adminPolicy},
		},
	}

	store.Policies = map[string]*Policy{
		adminPolicy: {
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
		},
	}
}
