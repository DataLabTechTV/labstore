package format

import (
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Size(size int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d B", size)
}

func Date(date time.Time) string {
	return date.Format(types.ISO8601)
}
