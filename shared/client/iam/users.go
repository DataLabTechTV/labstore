package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/iam"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func (c *Client) CreateUser(userName string) (*types.CreateUserResponse, error) {
	resp, err := c.DoRequest(iam.OpCreateUser, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.CreateUserResponse
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

func (c *Client) GetUser(userName string) (*types.GetUserResponse, error) {
	resp, err := c.DoRequest(iam.OpGetUser, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.GetUserResponse
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

func (c *Client) DeleteUser(userName string) (*types.DeleteUserResponse, error) {
	resp, err := c.DoRequest(iam.OpDeleteUser, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DeleteUserResponse
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

func (c *Client) CreateAccessKey(userName string) (*types.CreateAccessKeyResponse, error) {
	resp, err := c.DoRequest(iam.OpCreateAccessKey, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.CreateAccessKeyResponse
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

func (c *Client) ListAccessKeys(userName string) (*types.ListAccessKeysResponse, error) {
	resp, err := c.DoRequest(iam.OpListAccessKeys, "UserName", userName)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.ListAccessKeysResponse
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

func (c *Client) DeleteAccessKey(userName, accessKeyID string) (*types.DeleteAccessKeyResponse, error) {
	resp, err := c.DoRequest(iam.OpDeleteAccessKey, "UserName", userName, "AccessKeyId", accessKeyID)
	if err != nil {
		return nil, err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	if resp.StatusCode == http.StatusOK {
		var res types.DeleteAccessKeyResponse
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
