package iam

import (
	"fmt"
)

type ArnType string

const (
	ArnUser   ArnType = "user"
	ArnGroup  ArnType = "group"
	ArnPolicy ArnType = "policy"
	ArnS3     ArnType = "s3"
)

const defaultAccountID = "000000000001"

func toArn(arnType ArnType, path string) string {
	if arnType == ArnS3 {
		return fmt.Sprintf("arn:aws:s3:::%s", path)
	}

	return fmt.Sprintf(
		"arn:aws:iam::%s:%s%s",
		defaultAccountID,
		arnType,
		path,
	)
}
