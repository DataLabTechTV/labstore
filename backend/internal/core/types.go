package core

import (
	"encoding/xml"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05Z"

type Object struct {
	Key          string
	LastModified Timestamp
	ETag         string
	Size         int64
	Owner        *Owner
}

type Owner struct {
	ID          string
	DisplayName string // deprecated, but we'll support it
}

type Timestamp time.Time

func (t Timestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	formatted := time.Time(t).Format(ISO8601)
	return e.EncodeElement(formatted, start)
}
