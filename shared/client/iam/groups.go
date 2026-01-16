package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/iam"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
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

	var iamErr errs.IAMError
	if err := helper.ReadXML(resp.Body, &iamErr); err != nil {
		return nil, err
	}
	iamErr.StatusCode = resp.StatusCode

	return nil, &iamErr
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

	var iamErr errs.IAMError
	if err := helper.ReadXML(resp.Body, &iamErr); err != nil {
		return nil, err
	}
	iamErr.StatusCode = resp.StatusCode

	return nil, &iamErr
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

	var iamErr errs.IAMError
	if err := helper.ReadXML(resp.Body, &iamErr); err != nil {
		return nil, err
	}
	iamErr.StatusCode = resp.StatusCode

	return nil, &iamErr
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

	var iamErr errs.IAMError
	if err := helper.ReadXML(resp.Body, &iamErr); err != nil {
		return nil, err
	}
	iamErr.StatusCode = resp.StatusCode

	return nil, &iamErr
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

	var iamErr errs.IAMError
	if err := helper.ReadXML(resp.Body, &iamErr); err != nil {
		return nil, err
	}
	iamErr.StatusCode = resp.StatusCode

	return nil, &iamErr
}
