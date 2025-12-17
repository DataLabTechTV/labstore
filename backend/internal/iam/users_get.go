package iam

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"
)

type GetUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetUserResponse"`
	GetUserResult    *GetUserResult
	ResponseMetadata *ResponseMetadata
}

type GetUserResult struct {
	User *UserResult
}

func (store *Store) GetUserByName(name string) (*User, error) {
	var user User

	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE name = $1`, name); err != nil {
		slog.Error("get user by name", "err", err)
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ArnUser, user.UserID)
	if err != nil {
		return nil, err
	}
	user.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		user.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	groups, err := store.getGroupsByUserID(user.UserID)
	if err != nil {
		return nil, err
	}
	user.GroupIDs = make([]string, len(groups))
	for i, group := range groups {
		user.GroupIDs[i] = group.GroupID
	}

	if user.AccessKeyID.Valid {
		store.Users[user.AccessKeyID.String] = &CachedUser{
			user:     &user,
			loadedAt: time.Now(),
		}
	}

	return &user, nil
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
