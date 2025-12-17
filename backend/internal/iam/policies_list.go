package iam

import (
	"encoding/xml"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
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

func ListAttachedUserPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.GetUserByName(userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
		return
	}

	members := make([]*AttachedPoliciesMember, len(user.PolicyIDs))
	for i, policyID := range user.PolicyIDs {
		policy, err := store.GetPolicyByID(policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &ListAttachedUserPoliciesResponse{
		ListAttachedUserPoliciesResult: &ListAttachedUserPoliciesResult{
			AttachedPolicies: &AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func ListAttachedGroupPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	group, err := store.GetGroupByName(groupName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeGroup), groupName))
		return
	}

	members := make([]*AttachedPoliciesMember, len(group.PolicyIDs))
	for i, policyID := range group.PolicyIDs {
		policy, err := store.GetPolicyByID(policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &ListAttachedGroupPoliciesResponse{
		ListAttachedGroupPoliciesResult: &ListAttachedGroupPoliciesResult{
			AttachedPolicies: &AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
