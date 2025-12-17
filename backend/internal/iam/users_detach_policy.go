package iam

import (
	"encoding/xml"
	"net/http"
)

type DetachUserPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DetachUserPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

func DetachUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
