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

	return cmd
}

func NewObjectsPutCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "put BUCKET PATH LOCAL_PATH",
		Short: "Put S3 object",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(handlerKeyCtx).(*handlers.S3Handler)

			bucket := args[0]
			key := args[1]
			localPath := args[2]

			debug := helper.Must(cmd.Flags().GetBool("debug"))

			handler.PutObject(bucket, key, localPath, debug)
		},
	}

	return cmd
}
