package iam

import "github.com/DataLabTechTV/labstore/backend/config"

var Users map[string]string
var Policies map[string]PolicyFunc

type PolicyFunc func(userID string, resourceID string) bool

func Load() {
	Users = map[string]string{
		config.Env.AdminAccessKey: config.Env.AdminSecretKey,
	}

	Policies = map[string]PolicyFunc{
		config.Env.AdminAccessKey: func(bucket, op string) bool {
			return true
		},
	}
}
