package iam

import (
	"strings"

	"github.com/gobwas/glob"
)

func Resource(bucket, key string) string {
	resource := strings.TrimSuffix(bucket, "/") + "/" + strings.TrimPrefix(key, "/")
	return resource
}

func matchResource(bucket, key string, stmtResources []string) bool {
	resource := Resource(bucket, key)

	for _, stmtResource := range stmtResources {
		g := glob.MustCompile(stmtResource)
		if g.Match(resource) {
			return true
		}
	}

	return false
}
