package iam

import "log"

const Any = "*"

var store *Store

func Init() {
	store = NewStore()

	var err error

	err = store.open()
	if err != nil {
		log.Fatal(err)
	}

	store.setupAdmin()

	err = store.ensureSchema()
	if err != nil {
		log.Fatal(err)
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
