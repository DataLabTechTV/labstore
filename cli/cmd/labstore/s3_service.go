package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/spf13/cobra"
)

func NewServiceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "service",
		Short: "Handle S3 service operations",
	}

	cmd.AddCommand(NewServiceListObjectsCmd())

	return cmd
}

func NewServiceListObjectsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-buckets",
		Short: "List S3 buckets",
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			handler.ListBuckets()
		},
	}

	return cmd
}
