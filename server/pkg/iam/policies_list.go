package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
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

	members := make([]*types.AttachedPoliciesMember, len(user.PolicyIDs))
	for i, policyID := range user.PolicyIDs {
		policy, err := store.GetPolicyByID(ctx, policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &types.AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &types.ListAttachedUserPoliciesResponse{
		ListAttachedUserPoliciesResult: &types.ListAttachedUserPoliciesResult{
			AttachedPolicies: &types.AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
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

	members := make([]*types.AttachedPoliciesMember, len(group.PolicyIDs))
	for i, policyID := range group.PolicyIDs {
		policy, err := store.GetPolicyByID(ctx, policyID)
		if err != nil {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		members[i] = &types.AttachedPoliciesMember{
			PolicyName: policy.Name,
			PolicyArn:  policy.Arn,
		}
	}

	response := &types.ListAttachedGroupPoliciesResponse{
		ListAttachedGroupPoliciesResult: &types.ListAttachedGroupPoliciesResult{
			AttachedPolicies: &types.AttachedPolicies{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
