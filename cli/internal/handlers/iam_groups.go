package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/IllumiKnowLabs/labstore/cli/internal/errs"
	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
)

func (h *IAMHandler) CreateGroup(groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("CreateGroup: %s", groupName)))

	res, err := h.Client.CreateGroup(groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := groupMetadata(res.CreateGroupResult.Group)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) AddUserToGroup(userName, groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("AddUserToGroup: %s → %s", userName, groupName)))

	_, err := h.Client.AddUserToGroup(userName, groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "User added to group"))
	return nil
}

func (h *IAMHandler) GetGroup(groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("GetGroup: %s", groupName)))

	res, err := h.Client.GetGroup(groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := groupMetadata(res.GetGroupResult.Group)
	fmt.Println(metadata.Render())
	return nil
}

func (h *IAMHandler) DeleteGroup(groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("DeleteGroup: %s", groupName)))

	_, err := h.Client.DeleteGroup(groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "Group deleted"))
	return nil
}

func (h *IAMHandler) RemoveUserFromGroup(userName, groupName string) error {
	fmt.Println(render.Title(fmt.Sprintf("RemoveUserFromGroup: %s ↛  %s", userName, groupName)))

	_, err := h.Client.RemoveUserFromGroup(userName, groupName)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(http.StatusOK, "User removed from group"))
	return nil
}

func groupMetadata(res *types.GroupResult) render.Metadata {
	return render.Metadata{
		"Group ID":   render.NewString(res.GroupId),
		"Group Name": render.NewString(res.GroupName),
		"ARN":        render.NewString(res.Arn),
		"Path":       render.NewString(res.Path),
	}
}
