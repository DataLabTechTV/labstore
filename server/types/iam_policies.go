package types

import (
	"encoding/xml"
	"time"
)

type PolicyResult struct {
	XMLName          xml.Name `xml:"Policy"`
	PolicyName       string
	DefaultVersionId string
	PolicyId         string
	Path             string
	Arn              string
	AttachmentCount  int
	CreateDate       time.Time
	UpdateDate       time.Time
}

type CreatePolicyResponse struct {
	XMLName            xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreatePolicyResponse"`
	CreatePolicyResult *CreatePolicyResult
	ResponseMetadata   *ResponseMetadata
}

type CreatePolicyResult struct {
	Policy *PolicyResult
}

type AttachUserPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ AttachUserPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

type GetPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetPolicyResponse"`
	GetPolicyResult  *GetPolicyResult
	ResponseMetadata *ResponseMetadata
}

type GetPolicyResult struct {
	Policy *PolicyResult
}

type ListAttachedUserPoliciesResponse struct {
	XMLName                        xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAttachedUserPoliciesResponse"`
	ListAttachedUserPoliciesResult *ListAttachedUserPoliciesResult
	ResponseMetadata               *ResponseMetadata
}

type ListAttachedUserPoliciesResult struct {
	AttachedPolicies *AttachedPolicies
	IsTruncated      bool
}

type ListAttachedGroupPoliciesResponse struct {
	XMLName                         xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAttachedGroupPoliciesResponse"`
	ListAttachedGroupPoliciesResult *ListAttachedGroupPoliciesResult
	ResponseMetadata                *ResponseMetadata
}

type ListAttachedGroupPoliciesResult struct {
	AttachedPolicies *AttachedPolicies
	IsTruncated      bool
}

type AttachedPolicies struct {
	Member []*AttachedPoliciesMember `xml:"member"`
}

type AttachedPoliciesMember struct {
	PolicyName string
	PolicyArn  string
}

type DeletePolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeletePolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

type DetachUserPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DetachUserPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

type DetachGroupPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DetachGroupPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}
