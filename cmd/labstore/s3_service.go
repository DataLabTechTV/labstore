package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/handlers"
	"github.com/spf13/cobra"
)

func NewServiceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "service",
		Short: "Handle S3 service operations",
	}

	cmd.AddCommand(NewServiceListBucketsCmd())

	return cmd
}

func NewServiceListBucketsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-buckets",
		Short: "List S3 buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.ListBuckets()
		},
	}

	return cmd
}
