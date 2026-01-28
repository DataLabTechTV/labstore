package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/pkg/handlers"
	"github.com/spf13/cobra"
)

func NewPoliciesCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "policies",
		Short: "Handle IAM policy operations",
	}

	cmd.AddCommand(NewPoliciesCreateCmd())
	cmd.AddCommand(NewPoliciesAttachToUserCmd())
	cmd.AddCommand(NewPoliciesAttachToGroupCmd())
	cmd.AddCommand(NewPoliciesGetCmd())
	cmd.AddCommand(NewPoliciesListAttachedToUserCmd())
	cmd.AddCommand(NewPoliciesListAttachedToGroupCmd())
	cmd.AddCommand(NewPoliciesDeleteCmd())
	cmd.AddCommand(NewPoliciesDetachFromUserCmd())
	cmd.AddCommand(NewPoliciesDetachFromGroupCmd())

	return cmd
}

func NewPoliciesCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create POLICY_NAME POLICY_DOCUMENT_PATH",
		Short: "Create an IAM policy",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.CreatePolicy(args[0], args[1])
		},
	}

	return cmd
}

func NewPoliciesAttachToUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "attach-user USERNAME POLICY_ARN",
		Short: "Attach an IAM policy to an IAM user",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.AttachUserPolicy(args[0], args[1])
		},
	}

	return cmd
}

func NewPoliciesAttachToGroupCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "attach-group GROUPNAME POLICY_ARN",
		Short: "Attach an IAM policy to an IAM group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.AttachGroupPolicy(args[0], args[1])
		},
	}

	return cmd
}

func NewPoliciesGetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get POLICY_ARN",
		Short: "Display IAM policy information",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.GetPolicy(args[0])
		},
	}

	return cmd
}

func NewPoliciesListAttachedToUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-user-attached USERNAME",
		Short: "List IAM policies attached to an IAM user",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.ListAttachedUserPolicies(args[0])
		},
	}

	return cmd
}

func NewPoliciesListAttachedToGroupCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-group-attached GROUPNAME",
		Short: "List IAM policies attached to an IAM group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.ListAttachedGroupPolicies(args[0])
		},
	}

	return cmd
}

func NewPoliciesDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete POLICY_ARN",
		Short: "Delete IAM policy",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DeletePolicy(args[0])
		},
	}

	return cmd
}

func NewPoliciesDetachFromUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "detach-user USERNAME POLICY_ARN",
		Short: "Detach an IAM policy from an IAM user",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DetachUserPolicy(args[0], args[1])
		},
	}

	return cmd
}

func NewPoliciesDetachFromGroupCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "detach-group GROUPNAME POLICY_ARN",
		Short: "Detach an IAM policy from an IAM group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.IAMHandler)
			return handler.DetachGroupPolicy(args[0], args[1])
		},
	}

	return cmd
}
