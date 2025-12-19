package iam

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

const testUserName = "integration_test_user"

func TestCreateAccessKeyIntegration(t *testing.T) {
	helper.CheckFatal(os.Chdir("../../.."))

	config.Load(nil)
	Init()

	ctx := context.Background()

	_, err := store.CreateUser(ctx, testUserName)
	if err != nil {
		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			slog.Warn("user exists", "err", err)
		} else {
			t.Error(err)
		}
	}

	user, err := store.GetUserByName(ctx, testUserName)
	if err != nil {
		t.Error(err)
	}

	secretKey, err := store.CreateAccessKey(ctx, user)
	if err != nil {
		t.Error(err)
	}

	delete(store.CachedUsers, user.Name)
	fetchedUser, err := store.GetUserByAccessKey(ctx, user.AccessKeyID.String)
	if err != nil {
		t.Fatal(err)
	}

	if fetchedUser.UserID != user.UserID {
		t.Fatalf("expected %v, got %v", user.UserID, fetchedUser.UserID)
	}

	fetchedUserSecretKey, err := security.DecryptAESGCM(fetchedUser.SecretKey, config.Storage.MasterKeyPath)
	if err != nil {
		t.Error("failed to decrypt secret key after fetching user")
	}

	if fetchedUserSecretKey != secretKey {
		t.Error("decrypted secret key does not match set value")
	}
}
