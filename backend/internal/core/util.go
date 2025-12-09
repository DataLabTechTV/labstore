package core

import (
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
)

func BucketExists(bucket string) bool {
	path := filepath.Join(config.Server.Storage.Path, bucket)
	exists := helper.FileExists(path)
	return exists
}

func BucketKeyPath(bucket, key string) string {
	path := filepath.Join(config.Server.Storage.Path, bucket, key)
	return path
}
