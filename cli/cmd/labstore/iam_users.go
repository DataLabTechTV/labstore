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
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.CreateUser(args[0])
		},
	}

	return cmd
}

func NewUsersGetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get USERNAME",
		Short: "Display IAM user information",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.GetUser(args[0])
		},
	}

	return cmd
}

func NewUsersDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete USERNAME",
		Short: "Delete IAM user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.DeleteUser(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create USERNAME",
		Short: "Create IAM user access key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.CreateAccessKey(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysListCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list USERNAME",
		Short: "List IAM user access keys",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.ListAccessKeys(args[0])
		},
	}

	return cmd
}

func NewUsersAccessKeysDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete USERNAME",
		Short: "Delete IAM user access key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			handler.DeleteAccessKey(args[0])
		},
	}

	return cmd
}
