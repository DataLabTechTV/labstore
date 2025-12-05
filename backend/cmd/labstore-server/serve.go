package main

import (
	"log/slog"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/router"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Run backend server",
		Run: func(cmd *cobra.Command, args []string) {
			adminPass := cmd.Flags().Lookup("admin-pass")
			if adminPass != nil && adminPass.Changed {
				slog.Warn("setting admin pass via the command line is insecure")
			}

			router.Start()
		},
	}

	cmd.Flags().String("host", config.DefaultHost, "Listening hostname")
	cmd.Flags().Uint16("port", config.DefaultPort, "Listening port")
	cmd.Flags().String("storage-path", config.DefaultStoragePath, "Storage path for objects and internals")
	cmd.Flags().String("admin-user", config.DefaultAdminAccessKey, "Admin username / access key")
	cmd.Flags().String("admin-pass", config.DefaultAdminSecretKey, "Admin password / secret key")

	return cmd
}
