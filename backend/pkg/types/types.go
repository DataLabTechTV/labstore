package types

import (
	"encoding/xml"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05Z"

type ResponseMetadata struct {
	RequestId string
}

type Owner struct {
	ID          string
	DisplayName string // deprecated, but we'll support it
}

type Timestamp time.Time

func (t Timestamp) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	formatted := time.Time(t).Format(ISO8601)
	return enc.EncodeElement(formatted, start)
}
