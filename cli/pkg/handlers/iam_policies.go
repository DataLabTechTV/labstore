package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/cli/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/cli/pkg/render"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (h *IAMHandler) CreatePolicy(policyName, policyDocumentPath string) error {
	fmt.Println(render.Title(fmt.Sprintf("CreatePolicy: %s", policyName)))

	docFile, err := os.Open(policyDocumentPath)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}
	defer helper.CloseWithErr(docFile, &err)

	res, err := h.Client.CreatePolicy(policyName, docFile)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := policyMetadata(res.CreatePolicyResult.Policy)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) AttachUserPolicy(userName, policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("AttachUserPolicy: %s → %s", userName, policyArn)))

	_, err := h.Client.AttachUserPolicy(userName, policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Policy attached to user"))
	return nil
}

func (h *IAMHandler) AttachGroupPolicy(groupName, policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("AttachGroupPolicy: %s → %s", groupName, policyArn)))

	_, err := h.Client.AttachGroupPolicy(groupName, policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Policy attached to group"))
	return nil
}

func (h *IAMHandler) GetPolicy(policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("GetPolicy: %s", policyArn)))

	res, err := h.Client.GetPolicy(policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := policyMetadata(res.GetPolicyResult.Policy)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) ListAttachedUserPolicies(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("ListAttachedUserPolicies: %s", userName)))

	res, err := h.Client.ListAttachedUserPolicies(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	policiesMetadata := attachedPoliciesMetadata(res.ListAttachedUserPoliciesResult.AttachedPolicies.Member)
	for i, metadata := range policiesMetadata {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(metadata.Render())
	}
	return nil
}

func (h *IAMHandler) ListAttachedGroupPolicies(groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("ListAttachedGroupPolicies: %s", groupName)))

	res, err := h.Client.ListAttachedGroupPolicies(groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	policiesMetadata := attachedPoliciesMetadata(res.ListAttachedGroupPoliciesResult.AttachedPolicies.Member)
	for i, metadata := range policiesMetadata {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(metadata.Render())
	}
	return nil
}

func (h *IAMHandler) DeletePolicy(policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("DeletePolicy: %s", policyArn)))

	_, err := h.Client.DeletePolicy(policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Policy deleted"))
	return nil
}

func (h *IAMHandler) DetachUserPolicy(userName, policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("DetachUserPolicy: %s ↛  %s", userName, policyArn)))

	_, err := h.Client.DetachUserPolicy(userName, policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Policy detached from user"))
	return nil
}

func (h *IAMHandler) DetachGroupPolicy(groupName, policyArn string) error {
	fmt.Println(render.Title(fmt.Sprintf("DetachGroupPolicy: %s ↛  %s", groupName, policyArn)))

	_, err := h.Client.DetachGroupPolicy(groupName, policyArn)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Policy detached from group"))
	return nil
}

func policyMetadata(res *types.PolicyResult) render.Metadata {
	return render.Metadata{
		"PolicyName":       render.NewString(res.PolicyName),
		"DefaultVersionId": render.NewString(res.DefaultVersionId),
		"PolicyId":         render.NewString(res.PolicyId),
		"Path":             render.NewString(res.Path),
		"Arn":              render.NewString(res.Arn),
		"AttachmentCount":  render.NewNumber(res.AttachmentCount),
		"CreateDate":       render.NewDate(res.CreateDate),
		"UpdateDate":       render.NewDate(res.UpdateDate),
	}
}

func attachedPoliciesMetadata(res []*types.AttachedPoliciesMember) []render.Metadata {
	var list []render.Metadata

	for _, member := range res {
		metadata := render.Metadata{
			"Policy ARN":  render.NewString(member.PolicyName),
			"Policy Name": render.NewString(string(member.PolicyArn)),
		}
		list = append(list, metadata)
	}

	return list
}
