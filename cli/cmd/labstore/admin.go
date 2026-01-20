package main

import (
	"github.com/IllumiKnowLabs/labstore/client/admin"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
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

	cmd.Flags().String("admin-address-host", config.DefaultAdminAddressHost, "Listening host for admin server")
	cmd.Flags().String("admin-address-port", config.DefaultAdminAddressHost, "Listening port for admin server")
	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	return cmd
}
