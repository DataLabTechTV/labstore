package main

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/spf13/cobra"
)

func NewS3Cmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "s3",
		Short: "S3 client",
		Run: func(cmd *cobra.Command, args []string) {
			s3.Init()
		},
	}

	// TODO: need a way to manage credentials for other users (local config)
	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	cmd.Flags().String("s3-server-host", config.DefaultS3ServerHost, "Listening host for S3-compatible server")
	cmd.Flags().Uint16("s3-server-port", config.DefaultS3ServerPort, "Listening port for S3-compatible server")

	return cmd
}
