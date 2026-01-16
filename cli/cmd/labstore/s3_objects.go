package main

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/spf13/cobra"
)

func NewObjectsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "objects",
		Short: "Handle S3 object operations",
	}

	cmd.AddCommand(NewObjectsPutCmd())
	cmd.AddCommand(NewObjectsHeadCmd())
	cmd.AddCommand(NewObjectsGetCmd())
	cmd.AddCommand(NewObjectsDeleteCmd())

	return cmd
}

func NewObjectsPutCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "put BUCKET PATH LOCAL_PATH",
		Short: "Put S3 object",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)

			bucket := args[0]
			key := args[1]
			localPath := args[2]

			debug := helper.Must(cmd.Flags().GetBool("debug"))

			return handler.PutObject(bucket, key, localPath, debug)
		},
	}

	return cmd
}

func NewObjectsHeadCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "head BUCKET KEY",
		Short: "Metadata for an S3 object",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.HeadObject(args[0], args[1])
		},
	}

	return cmd
}

func NewObjectsGetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get BUCKET PATH LOCAL_PATH",
		Short: "Download an S3 object",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)

			bucket := args[0]
			key := args[1]
			localPath := args[2]

			debug := helper.Must(cmd.Flags().GetBool("debug"))

			return handler.GetObject(bucket, key, localPath, debug)
		},
	}

	return cmd
}

func NewObjectsDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete BUCKET KEYS...",
		Short: "Delete one or multiple S3 objects",
		Long:  "Delete one or multiple S3 objects, either using DeleteObject (DELETE) or DeleteObjects (POST ?delete)",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)
			return handler.DeleteObjects(args[0], args[1:]...)
		},
	}

	return cmd
}
