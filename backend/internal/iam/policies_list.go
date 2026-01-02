package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func ListAttachedUserPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
		return
	}

	members := make([]*t.AttachedPoliciesMember, len(user.PolicyIDs))
	for i, policyID := range user.PolicyIDs {
		policy, err := store.GetPolicyByID(ctx, policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &t.AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &t.ListAttachedUserPoliciesResponse{
		ListAttachedUserPoliciesResult: &t.ListAttachedUserPoliciesResult{
			AttachedPolicies: &t.AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &t.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func ListAttachedGroupPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	ctx := r.Context()

	group, err := store.GetGroupByName(ctx, groupName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeGroup), groupName))
		return
	}

	members := make([]*t.AttachedPoliciesMember, len(group.PolicyIDs))
	for i, policyID := range group.PolicyIDs {
		policy, err := store.GetPolicyByID(ctx, policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &t.AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &t.ListAttachedGroupPoliciesResponse{
		ListAttachedGroupPoliciesResult: &t.ListAttachedGroupPoliciesResult{
			AttachedPolicies: &t.AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &t.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
