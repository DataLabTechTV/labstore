package iam

import "github.com/DataLabTechTV/labstore/config"

var Users map[string]string
var Policies map[string]func(string, string) bool

func Load() {
	Users = map[string]string{
		config.Env.AdminAccessKey: config.Env.AdminSecretKey,
	}

	Policies = map[string]func(string, string) bool{
		config.Env.AdminAccessKey: func(bucket, op string) bool {
			return true
		},
	}
}
