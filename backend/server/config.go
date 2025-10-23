package server

import (
	"os"

	"github.com/joho/godotenv"
)

var storageRoot string
var users map[string]string
var userPolicies map[string]func(string, string) bool

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		os.Exit(2)
	}

	storageRoot = os.Getenv("LS_STORAGE_ROOT")

	var admin_access_key = os.Getenv("LS_ADMIN_ACCESS_KEY")
	var admin_secret_key = os.Getenv("LS_ADMIN_SECRET_KEY")

	users = map[string]string{
		admin_access_key: admin_secret_key,
	}

	userPolicies = map[string]func(string, string) bool{
		admin_access_key: func(bucket, op string) bool {
			return true
		},
	}
}
