package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/spf13/cobra"
)

func NewUsersCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "users",
		Short: "Handle IAM user operations",
	}

	cmd.AddCommand(NewUsersAccessKeysCmd())
	cmd.AddCommand(NewUsersCreateCmd())
	cmd.AddCommand(NewUsersGetCmd())
	cmd.AddCommand(NewUsersDeleteCmd())

	return cmd
}

func NewUsersAccessKeysCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "access-keys",
		Short: "Handle IAM access key operations",
	}

	cmd.AddCommand(NewUsersAccessKeysCreateCmd())
	cmd.AddCommand(NewUsersAccessKeysListCmd())
	cmd.AddCommand(NewUsersAccessKeysDeleteCmd())

	return cmd
}

func NewUsersCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create USERNAME",
		Short: "Create an IAM user",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.CreateUser(args[0])
		},
	}

	return cmd
}

func NewUsersGetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get USERNAME",
		Short: "Display IAM user information",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.GetUser(args[0])
		},
	}

	return cmd
}

func NewUsersDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete USERNAME",
		Short: "Delete IAM user",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DeleteUser(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create USERNAME",
		Short: "Create IAM user access key",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.CreateAccessKey(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysListCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list USERNAME",
		Short: "List IAM user access keys",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.ListAccessKeys(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete USERNAME ACCESS_KEY_ID",
		Short: "Delete IAM user access key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DeleteAccessKey(args[0], args[1])
		},
	}

	return cmd
}
