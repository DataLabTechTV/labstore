package iam

import (
	"encoding/json"

	"github.com/IllumiKnowLabs/labstore/server/core"
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
	path := core.ObjectPath(bucket, key)
	resource := toArn(ArnS3, path)
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
