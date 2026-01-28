package helper

import (
	"encoding/xml"
	"io"

	"github.com/IllumiKnowLabs/labstore/server/helper"
)

func XMLReader(v any) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		defer helper.CloseWithErr(pw, nil)
		enc := xml.NewEncoder(pw)
		if err := enc.Encode(v); err != nil {
			pw.CloseWithError(err)
		}
	}()

	return pr
}
