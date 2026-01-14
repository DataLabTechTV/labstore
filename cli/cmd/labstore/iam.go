package main

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/client/iam"
	"github.com/spf13/cobra"
)

func NewIAMCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "iam",
		Short: "IAM client, designed for learning",
		Run: func(cmd *cobra.Command, args []string) {
			iam.Init()
		},
	}

	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	cmd.Flags().String("iam-server-host", config.DefaultIAMServerHost, "Listening host for IAM server")
	cmd.Flags().String("iam-server-port", config.DefaultIAMServerHost, "Listening port for IAM server")

	return cmd
}
