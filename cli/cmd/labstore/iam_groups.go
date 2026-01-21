package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/spf13/cobra"
)

func NewGroupsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "groups",
		Short: "Handle IAM group operations",
	}

	cmd.AddCommand(NewGroupsCreateCmd())
	cmd.AddCommand(NewGroupsAddUserCmd())
	cmd.AddCommand(NewGroupsGetCmd())
	cmd.AddCommand(NewGroupsDeleteCmd())
	cmd.AddCommand(NewGroupsRemoveUserCmd())

	return cmd
}

func NewGroupsCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create GROUPNAME",
		Short: "Create an IAM group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.CreateGroup(args[0])
		},
	}

	return cmd
}

func NewGroupsAddUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add-user USERNAME GROUPNAME",
		Short: "Add an IAM user to an IAM group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.AddUserToGroup(args[0], args[1])
		},
	}

	return cmd
}

func NewGroupsGetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get GROUPNAME",
		Short: "Display IAM group information",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.GetGroup(args[0])
		},
	}

	return cmd
}

func NewGroupsDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete GROUPNAME",
		Short: "Delete IAM group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DeleteGroup(args[0])
		},
	}

	return cmd
}

func NewGroupsRemoveUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remove-user USERNAME GROUPNAME",
		Short: "Remove an IAM user from an IAM group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.RemoveUserFromGroup(args[0], args[1])
		},
	}

	return cmd
}
