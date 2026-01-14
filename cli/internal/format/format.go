package format

import (
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Size(size int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d B", size)
}

func Date(date types.Timestamp) string {
	return fmt.Sprintf("[%s]", time.Time(date).Format(types.ISO8601))
}
