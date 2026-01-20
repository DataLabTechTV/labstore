package client

import (
	"encoding/xml"
	"io"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
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
