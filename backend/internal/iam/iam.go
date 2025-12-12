package iam

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

const Any = "*"

const masterKeyFilename = "master.key"

var store *Store

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
	user, ok := store.Users[accessKey]
	if !ok {
		return false
	}

	allowed := false

	for _, policyID := range user.PolicyIDs {
		policy, ok := store.Policies[policyID]
		if !ok {
			continue
		}

		for _, stmt := range policy.Document.Statement {
			if matchAction(action, stmt.Actions) && matchResource(bucket, key, stmt.Resources) {
				if stmt.Effect == deny {
					return false
				}

				if stmt.Effect == allow {
					allowed = true
				}
			}
		}
	}

	return allowed
}

func ensureDirectories() error {
	slog.Debug("ensuring iam directories")

	if err := os.MkdirAll(config.Storage.MetadataPath, 0750); err != nil {
		return err
	}

	if err := os.MkdirAll(config.Storage.KeysDir, 0600); err != nil {
		return err
	}

	return nil
}

func ensureMasterKey() error {
	key, err := security.GenerateAES256MasterKey()
	if err != nil {
		return err
	}

	keyPath := filepath.Join(config.Storage.KeysDir, masterKeyFilename)

	if _, err := os.Stat(keyPath); os.IsExist(err) {
		return nil
	}

	if err := os.WriteFile(keyPath, key, 0400); err != nil {
		return err
	}

	return nil
}

func initStore() error {
	store = NewStore()

	if err := store.open(); err != nil {
		return err
	}

	store.setupAdmin()

	if err := store.ensureSchema(); err != nil {
		return nil
	}

	return nil
}
