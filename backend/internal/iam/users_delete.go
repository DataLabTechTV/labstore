package iam

import (
	"encoding/xml"
	"net/http"
)

type DeleteUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteUserResponse"`
	ResponseMetadata *ResponseMetadata
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
