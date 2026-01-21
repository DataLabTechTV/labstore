package types

import (
	"encoding/xml"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05Z"
const CompactISO8601 = "20060102T150405Z"

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

func (t *Timestamp) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return err
	}

	parsed, err := time.Parse(ISO8601, v)
	if err != nil {
		return err
	}

	*t = Timestamp(parsed)
	return nil
}
