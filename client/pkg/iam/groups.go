package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/iam"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (c *Client) CreateGroup(groupName string) (*types.CreateGroupResponse, error) {
	resp, err := c.DoRequest(iam.OpCreateGroup, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.CreateGroupResponse
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

func (c *Client) AddUserToGroup(userName, groupName string) (*types.AddUserToGroupResponse, error) {
	resp, err := c.DoRequest(iam.OpAddUserToGroup, "UserName", userName, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.AddUserToGroupResponse
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

func (c *Client) GetGroup(groupName string) (*types.GetGroupResponse, error) {
	resp, err := c.DoRequest(iam.OpGetGroup, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.GetGroupResponse
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

func (c *Client) DeleteGroup(groupName string) (*types.DeleteGroupResponse, error) {
	resp, err := c.DoRequest(iam.OpDeleteGroup, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DeleteGroupResponse
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

func (c *Client) RemoveUserFromGroup(userName, groupName string) (*types.RemoveUserFromGroupResponse, error) {
	resp, err := c.DoRequest(iam.OpRemoveUserFromGroup, "UserName", userName, "GroupName", groupName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.RemoveUserFromGroupResponse
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
