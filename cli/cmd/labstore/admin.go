package main

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/client/admin"
	"github.com/spf13/cobra"
)

func NewAdminCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "admin",
		Short: "Admin client",
		Run: func(cmd *cobra.Command, args []string) {
			admin.Init()
		},
	}

	cmd.Flags().String("admin-server-host", config.DefaultAdminServerHost, "Listening host for admin server")
	cmd.Flags().String("admin-server-port", config.DefaultAdminServerHost, "Listening port for admin server")
	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	return cmd
}
