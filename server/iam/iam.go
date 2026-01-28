package iam

import (
	"context"
	"log/slog"
	"os"

	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/IllumiKnowLabs/labstore/server/helper"
	"github.com/IllumiKnowLabs/labstore/server/security"
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
	slog.Debug("check policy", "accessKey", accessKey, "bucket", bucket, "key", key, "action", action)

	ctx := context.Background()

	user, err := store.GetUserByAccessKey(ctx, accessKey)
	if err != nil {
		return false
	}

	var policyIDs []string
	uniquePolicyIDs := make(map[string]bool)

	for _, policyID := range user.PolicyIDs {
		if _, ok := uniquePolicyIDs[policyID]; ok {
			continue
		}
		uniquePolicyIDs[policyID] = true
		policyIDs = append(policyIDs, policyID)
	}

	for _, groupID := range user.GroupIDs {
		group, err := store.GetGroupByID(ctx, groupID)
		if err != nil {
			slog.Warn("group not found", "err", err)
		}

		for _, policyID := range group.PolicyIDs {
			if _, ok := uniquePolicyIDs[policyID]; ok {
				continue
			}
			uniquePolicyIDs[policyID] = true
			policyIDs = append(policyIDs, policyID)
		}
	}

	allowed := false

	for _, policyID := range policyIDs {
		policy, err := store.GetPolicyByID(ctx, policyID)
		if err != nil {
			slog.Warn("policy not found", "policy", policy, "err", err)
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

	return allowed
}

func ensureDirectories() error {
	slog.Debug("ensuring iam directories")

	if err := os.MkdirAll(config.App.Server.Storage.MetadataPath, 0o750); err != nil {
		return err
	}

	if err := os.MkdirAll(config.App.Server.Storage.KeysDir, 0o700); err != nil {
		return err
	}

	return nil
}

func ensureMasterKey() error {
	if helper.FileExists(config.App.Server.Storage.MasterKeyPath) {
		return nil
	}

	key, err := security.GenerateMasterKey()
	if err != nil {
		return err
	}

	if err := os.WriteFile(config.App.Server.Storage.MasterKeyPath, key, 0o600); err != nil {
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
