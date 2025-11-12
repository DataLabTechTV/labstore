package object

import "github.com/DataLabTechTV/labstore/backend/internal/core"

type Object struct {
	Key          string
	LastModified core.Timestamp
	ETag         string
	Size         int64
	Owner        *core.Owner
}
