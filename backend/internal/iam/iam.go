package iam

import (
	"log/slog"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

const Any = "*"

var store *Store

func GetStore() *Store {
	return store
}

func Init() {
	if err := ensureDirectories(); err != nil {
		slog.Error("could not create iam directories", "err", err)
		os.Exit(1)
	}

	if err := ensureMasterKey(); err != nil {
		slog.Error("could not access master key", "err", err)
		os.Exit(1)
	}

	if err := initStore(); err != nil {
		slog.Error("could not create initialize iam store", "err", err)
		os.Exit(1)
	}
}

func CheckPolicy(accessKey, bucket, key string, action Action) bool {
	user, err := store.GetUserByAccessKey(accessKey)
	if err != nil {
		return false
	}

	allowed := false

	for _, policyID := range user.PolicyIDs {
		policy, ok := store.Policies[policyID]
		if !ok {
			continue
		}

		for _, stmt := range policy.Document.Statement {
			if matchAction(action, stmt.Action) && matchResource(bucket, key, stmt.Resource) {
				if stmt.Effect == deny {
					return false
				}

				if stmt.Effect == allow {
					allowed = true
				}
			}
		}
	}

	// TODO: group policy check

	return allowed
}

func ensureDirectories() error {
	slog.Debug("ensuring iam directories")

	if err := os.MkdirAll(config.Storage.MetadataPath, 0750); err != nil {
		return err
	}

	if err := os.MkdirAll(config.Storage.KeysDir, 0700); err != nil {
		return err
	}

	return nil
}

func ensureMasterKey() error {
	if helper.FileExists(config.Storage.MasterKeyPath) {
		return nil
	}

	key, err := security.GenerateMasterKey()
	if err != nil {
		return err
	}

	if err := os.WriteFile(config.Storage.MasterKeyPath, key, 0600); err != nil {
		return err
	}

	return nil
}

func initStore() error {
	store = NewStore()

	if err := store.open(); err != nil {
		return err
	}

	if err := store.setupAdmin(); err != nil {
		return err
	}

	if err := store.ensureSchema(); err != nil {
		return err
	}

	return nil
}
