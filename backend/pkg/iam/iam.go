package iam

import "github.com/IllumiKnowLabs/labstore/backend/internal/config"

var Users map[string]string
var Policies map[string]PolicyFunc

type PolicyFunc func(userID string, resourceID string) bool

func Load() {
	Users = map[string]string{
		config.Config.Server.Admin.AccessKey: config.Config.Server.Admin.SecretKey,
	}

	Policies = map[string]PolicyFunc{
		config.Config.Server.Admin.AccessKey: func(bucket, op string) bool {
			return true
		},
	}
}

func CheckPolicy(accessKey, bucket, op string) bool {
	if polFunc, ok := Policies[accessKey]; ok {
		return polFunc(bucket, op)
	}
	return false
}
