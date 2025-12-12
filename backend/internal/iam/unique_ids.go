package iam

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UniqueIDPrefix string

const (
	IAMUserUniqueID       UniqueIDPrefix = "AIDA"
	ManagedPolicyUniqueID UniqueIDPrefix = "ANPA"
)

func GenerateUniqueID(prefix UniqueIDPrefix) string {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
	uniqueID := fmt.Sprintf("%s%s", prefix, id)
	return uniqueID
}
