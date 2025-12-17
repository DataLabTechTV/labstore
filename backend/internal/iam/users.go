package iam

import (
	"database/sql"
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
