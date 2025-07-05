package server

var storageRoot string
var users map[string]string
var userPolicies map[string]func(string, string) bool

func init() {
	storageRoot = "./data" // local dir to store buckets & objects

	// In-memory user with access key & secret key (hardcoded for MVP)
	users = map[string]string{
		"myaccesskey": "mysecretkey",
	}

	// Simple user policies (for MVP: allow all operations)
	userPolicies = map[string]func(string, string) bool{
		"myaccesskey": func(bucket, op string) bool {
			return true
		},
	}
}
