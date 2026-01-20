package iam

import (
	"io"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/iam"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (c *Client) CreatePolicy(policyName string, policyReader io.Reader) (*types.CreatePolicyResponse, error) {
	buf, err := io.ReadAll(policyReader)
	if err != nil {
		return nil, err
	}
	policyDocument := string(buf)

	resp, err := c.DoRequest(iam.OpCreatePolicy, "PolicyName", policyName, "PolicyDocument", policyDocument)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.CreatePolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) AttachUserPolicy(userName, policyArn string) (*types.AttachUserPolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpAttachUserPolicy, "UserName", userName, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.AttachUserPolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) AttachGroupPolicy(groupName, policyArn string) (*types.AttachUserPolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpAttachGroupPolicy, "GroupName", groupName, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.AttachUserPolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) GetPolicy(policyArn string) (*types.GetPolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpGetPolicy, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.GetPolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) ListAttachedUserPolicies(userName string) (*types.ListAttachedUserPoliciesResponse, error) {
	resp, err := c.DoRequest(iam.OpListAttachedUserPolicies, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.ListAttachedUserPoliciesResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) ListAttachedGroupPolicies(groupName string) (*types.ListAttachedGroupPoliciesResponse, error) {
	resp, err := c.DoRequest(iam.OpListAttachedGroupPolicies, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.ListAttachedGroupPoliciesResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) DeletePolicy(policyArn string) (*types.DeletePolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpDeletePolicy, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DeletePolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) DetachUserPolicy(userName, policyArn string) (*types.DetachUserPolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpDetachUserPolicy, "UserName", userName, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DetachUserPolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}

func (c *Client) DetachGroupPolicy(groupName, policyArn string) (*types.DetachGroupPolicyResponse, error) {
	resp, err := c.DoRequest(iam.OpDetachGroupPolicy, "GroupName", groupName, "PolicyArn", policyArn)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DetachGroupPolicyResponse
		if err := helper.ReadXML(resp.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	}

	var iamErrResp errs.IAMErrorResponse
	if err := helper.ReadXML(resp.Body, &iamErrResp); err != nil {
		return nil, err
	}
	iamErrResp.Error.StatusCode = resp.StatusCode

	return nil, iamErrResp.Error
}
