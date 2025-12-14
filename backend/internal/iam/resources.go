package iam

import (
	"encoding/json"
	"strings"

	"github.com/gobwas/glob"
)

type Resources []string

func (r *Resources) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*r = []string{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}
	*r = multi
	return nil
}

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
