package main

import (
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/spf13/cobra"
)

func NewBucketsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "buckets",
		Short: "Handle S3 bucket operations",
	}

	cmd.AddCommand(NewBucketsCreateCmd())
	cmd.AddCommand(NewBucketsHeadCmd())
	cmd.AddCommand(NewBucketsListObjectsCmd())
	cmd.AddCommand(NewBucketsDeleteCmd())

	return cmd
}

func NewBucketsCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create BUCKET",
		Short: "Create an S3 bucket",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.CreateBucket(args[0])
		},
	}

	return cmd
}

func NewBucketsHeadCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "head BUCKET",
		Short: "Status for an S3 bucket",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.HeadBucket(args[0])
		},
	}

	return cmd
}

func NewBucketsListObjectsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-objects BUCKET [PATH]",
		Short: "List S3 objects",
		Long:  "List S3 objects for the specified BUCKET, with an optional PATH globbing expression",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)

			bucket := args[0]

			var key *string
			if len(args) >= 2 {
				key = helper.StringPtr(args[1])
			} else {
				key = nil
			}

			return handler.ListObjects(bucket, key)
		},
	}

	return cmd
}

func NewBucketsDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete BUCKET",
		Short: "Delete S3 bucket",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.DeleteBucket(args[0])
		},
	}

	return cmd
}
