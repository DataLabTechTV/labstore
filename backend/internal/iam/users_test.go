package iam

import (
	"log/slog"
	"os"
	"testing"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

const testUserName = "integration_test_user"

func TestCreateAccessKeyIntegration(t *testing.T) {
	helper.CheckFatal(os.Chdir("../../.."))

	config.Load(nil)
	Init()

	_, iamErr := CreateUser(testUserName)
	if iamErr != nil {
		if iamErr.Code == "EntityAlreadyExists" {
			slog.Warn("user exists", "err", iamErr)
		} else {
			t.Error(iamErr)
		}
	}

	user, err := GetUserByName(testUserName)
	if err != nil {
		t.Error(err)
	}

	_, err = CreateAccessKey(user)
	if err != nil {
		t.Error(err)
	}

	fetchedUser, err := GetUserByAccessKey(user.AccessKeyID.String)
	if err != nil {
		t.Fatal(err)
	}

	if fetchedUser.UserID != user.UserID {
		t.Fatalf("expected %v, got %v", user.UserID, fetchedUser.UserID)
	}
}
