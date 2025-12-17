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

type AttachedPolicies struct {
	Member *AttachedPoliciesMember `xml:"member"`
}

type AttachedPoliciesMember struct {
	PolicyName string
	PolicyArn  string
}

func ListAttachedUserPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func ListAttachedGroupPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
