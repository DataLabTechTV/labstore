package core

import (
	"path/filepath"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
)

func BucketExists(bucket string) bool {
	path := BucketSystemPath(bucket)
	exists := helper.FileExists(path)
	return exists
}

func BucketSystemPath(bucket string) string {
	path := filepath.Join(config.Storage.ObjectsPath, bucket)
	return path
}

func ObjectSystemPath(bucket, key string) string {
	path := filepath.Join(config.Storage.ObjectsPath, bucket, key)
	return path
}

func ObjectPath(bucket, key string) string {
	path := strings.TrimSuffix(bucket, "/") + "/" + strings.TrimPrefix(key, "/")
	return path
}
