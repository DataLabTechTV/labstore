package types

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestTimestampMarshalXML(t *testing.T) {
	ts := Timestamp(time.Date(2025, time.November, 12, 18, 40, 0, 0, time.UTC))

	const expected = "<TestTimestampMarshalXML><LastModified>2025-11-12T18:40:00Z</LastModified></TestTimestampMarshalXML>"

	type Data struct {
		XMLName      xml.Name  `xml:"TestTimestampMarshalXML"`
		LastModified Timestamp `xml:"LastModified"`
	}

	data := Data{
		LastModified: ts,
	}

	var b strings.Builder

	if err := xml.NewEncoder(&b).Encode(data); err != nil {
		t.Fatal(err)
	}

	if b.String() != expected {
		t.Error("Encoding doesn't match expected XML")
	}
}
