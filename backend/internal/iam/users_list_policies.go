package iam

import (
	"encoding/xml"
	"net/http"
)

type ListAttachedUserPoliciesResponse struct {
	XMLName                        xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAttachedUserPoliciesResponse"`
	ListAttachedUserPoliciesResult *ListAttachedUserPoliciesResult
	ResponseMetadata               *ResponseMetadata
}

type ListAttachedUserPoliciesResult struct {
	AttachedPolicies *AttachedPolicies
	IsTruncated      bool
}

//nolint:unused
type AttachedPolicies struct {
	member *AttachedPoliciesMember
}

type AttachedPoliciesMember struct {
	PolicyName string
	PolicyArn  string
}

func ListAttachedUserPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
