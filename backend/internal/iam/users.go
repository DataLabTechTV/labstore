package iam

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

const (
	defaultAdminUserName = "Administrator"
	defaultUserPath      = "/"
)

type CachedUser struct {
	User        *User
	LoadedAt    time.Time
	NeverExpire bool
}

type User struct {
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Arn    string `db:"arn"`

	AccessKeyID sql.NullString `db:"access_key"`
	SecretKey   []byte         `db:"secret_key"`
	Salt        []byte         `db:"salt"`

	GroupIDs  []string
	PolicyIDs []string
}

type UserResult struct {
	Path     string
	UserName string
	UserId   string
	Arn      string
}

func (user *User) Result() *UserResult {
	userPath := "/"

	return &UserResult{
		Path:     userPath,
		UserName: user.Name,
		UserId:   user.UserID,
		Arn:      user.Arn,
	}
}

func (user *User) EncryptedData() *security.EncryptedData {
	return &security.EncryptedData{
		Value: user.SecretKey,
		Salt:  user.Salt,
	}
}

func (store *Store) GetUserByAccessKey(accessKey string) (*User, error) {
	if cachedUser, ok := store.CachedUsers[accessKey]; ok {
		if cachedUser.NeverExpire || time.Since(cachedUser.LoadedAt) < store.ttl {
			return cachedUser.User, nil
		}

		slog.Debug("invalidating cached user", "accessKey", accessKey)
		delete(store.CachedUsers, accessKey)
	}

	var user User
	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE access_key = $1`, accessKey); err != nil {
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

	store.CachedUsers[user.AccessKeyID.String] = &CachedUser{
		User:     &user,
		LoadedAt: time.Now(),
	}

	return &user, nil
}

func (store *Store) getUsersByGroupID(groupID string) ([]*User, error) {
	var users []*User

	query := `
	SELECT * FROM users WHERE user_id = (
		SELECT user_id FROM group_users WHERE group_id = $1
	)
	`

	if err := store.readDB.Select(&users, query, groupID); err != nil {
		slog.Error("get groups by user id", "err", err)
		return nil, err
	}

	return users, nil
}

func (store *Store) setupAdmin() error {
	encryptedData, err := security.EncryptAESGCM(config.Admin.Auth.SecretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return err
	}

	store.CachedUsers[config.Admin.Auth.AccessKey] = &CachedUser{
		User: &User{
			UserID: GenerateUniqueID(IAMUserUniqueID),
			Name:   defaultAdminUserName,
			Arn:    toArn(ArnUser, defaultUserPath+defaultAdminUserName),

			AccessKeyID: sql.NullString{String: config.Admin.Auth.AccessKey, Valid: true},
			SecretKey:   encryptedData.Value,
			Salt:        encryptedData.Salt,

			GroupIDs:  []string{},
			PolicyIDs: []string{adminPolicy},
		},
		LoadedAt:    time.Now(),
		NeverExpire: true,
	}

	store.CachedPolicies[adminPolicy] = &CachedPolicy{
		Policy: &Policy{
			PolicyID: adminPolicy,
			Document: &PolicyDocument{
				Version: latestPolicyDocumentVersion,
				Statement: []Statement{
					{
						Effect:   allow,
						Action:   []Action{Action(Any)},
						Resource: []string{Any},
					},
				},
			},
		},
		LoadedAt:    time.Now(),
		NeverExpire: true,
	}

	return nil
}
