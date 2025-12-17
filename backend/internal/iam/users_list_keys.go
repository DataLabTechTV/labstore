package iam

import (
	"encoding/xml"
	"net/http"
)

type ListAccessKeysResponse struct {
	XMLName              xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAccessKeysResponse"`
	ListAccessKeysResult *ListAccessKeysResult
	ResponseMetadata     *ResponseMetadata
}

type ListAccessKeysResult struct {
	UserName          string
	AccessKeyMetadata *AccessKeyMetadata
	IsTruncated       bool
}

//nolint:unused
type AccessKeyMetadata struct {
	member []*AccessKeyMetadataMember
}

type AccessKeyMetadataMember struct {
	UserName    string
	AccessKeyId string
	Status      AccessKeyStatus
}

func ListAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
