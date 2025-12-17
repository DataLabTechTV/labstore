package iam

import (
	"encoding/xml"
	"net/http"
)

type DeleteAccessKeyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteAccessKeyResponse"`
	ResponseMetadata *ResponseMetadata
}

func DeleteAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
