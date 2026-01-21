package iam

import (
	"database/sql"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/security"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
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

	GroupIDs  []string
	PolicyIDs []string
}

func (user *User) Result() *types.UserResult {
	userPath := "/"

	return &types.UserResult{
		Path:     userPath,
		UserName: user.Name,
		UserId:   user.UserID,
		Arn:      user.Arn,
	}
}

func (store *Store) setupAdmin() error {
	encrypted, err := security.EncryptAESGCM(
		config.App.Server.Admin.Auth.SecretKey,
		config.App.Server.Storage.MasterKeyPath,
	)
	if err != nil {
		return err
	}

	store.CachedUsers[config.App.Server.Admin.Auth.AccessKey] = &CachedUser{
		User: &User{
			UserID: GenerateUniqueID(IAMUserUniqueID),
			Name:   defaultAdminUserName,
			Arn:    toArn(ArnUser, defaultUserPath+defaultAdminUserName),

			AccessKeyID: sql.NullString{String: config.App.Server.Admin.Auth.AccessKey, Valid: true},
			SecretKey:   encrypted,

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
