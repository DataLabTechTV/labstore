package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/IllumiKnowLabs/labstore/cli/internal/errs"
	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
)

func (h *IAMHandler) CreateUser(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("CreateUser: %s", userName)))

	res, err := h.Client.CreateUser(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := userMetadata(res.CreateUserResult.User)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) GetUser(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("GetUsers: %s", userName)))

	res, err := h.Client.GetUser(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := userMetadata(res.GetUserResult.User)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) DeleteUser(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("DeleteUser: %s", userName)))

	_, err := h.Client.DeleteUser(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "User deleted"))
	return nil
}

func (h *IAMHandler) CreateAccessKey(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("CreateAccessKey: %s", userName)))

	res, err := h.Client.CreateAccessKey(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := accessKeyMetadata(res.CreateAccessKeyResult.AccessKey)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) ListAccessKeys(userName string) error {
	fmt.Println(render.Title(fmt.Sprintf("ListAccessKeys: %s", userName)))

	res, err := h.Client.ListAccessKeys(userName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	keysMetadata := accessKeyMembersMetadata(res.ListAccessKeysResult.AccessKeyMetadata.Member)
	for i, metadata := range keysMetadata {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(metadata.Render())
	}

	return nil
}

func (h *IAMHandler) DeleteAccessKey(userName, accessKeyID string) error {
	fmt.Println(render.Title(fmt.Sprintf("DeleteAccessKey: %s", userName)))

	_, err := h.Client.DeleteAccessKey(userName, accessKeyID)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Access key deleted"))
	return nil
}

func userMetadata(res *types.UserResult) render.Metadata {
	return render.Metadata{
		"User ID":   render.NewString(res.UserId),
		"User Name": render.NewString(res.UserName),
		"ARN":       render.NewString(res.Arn),
		"Path":      render.NewString(res.Path),
	}
}

func accessKeyMetadata(res *types.AccessKeyResult) render.Metadata {
	return render.Metadata{
		"Access Key ID":     render.NewString(res.AccessKeyId),
		"Secret Access Key": render.NewString(res.SecretAccessKey),
		"Status":            render.NewString(string(res.Status)),
	}
}

func accessKeyMembersMetadata(res []*types.AccessKeyMetadataMember) []render.Metadata {
	var list []render.Metadata

	for _, member := range res {
		metadata := render.Metadata{
			"Access Key ID": render.NewString(member.AccessKeyId),
			"Status":        render.NewString(string(member.Status)),
		}
		list = append(list, metadata)
	}

	return list
}
