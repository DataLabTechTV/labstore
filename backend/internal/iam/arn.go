package iam

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/internal/constants"
)

const (
	ArnUser   ArnType = "user"
	ArnGroup  ArnType = "group"
	ArnPolicy ArnType = "policy"
)

const defaultAccountID = "000000000001"

type ArnType string

func toArn(arnType ArnType, path, description string) string {
	return fmt.Sprintf(
		"arn:%s:iam::%s:%s%s%s",
		strings.ToLower(constants.Name),
		defaultAccountID,
		arnType,
		path,
		description,
	)
}
