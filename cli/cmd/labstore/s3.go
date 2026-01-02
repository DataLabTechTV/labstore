package main

import (
	"log/slog"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/spf13/cobra"
)

func NewS3Cmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "s3",
		Short: "S3 client",
		Run: func(cmd *cobra.Command, args []string) {
			credentials.Init()

			profileName := helper.Must(cmd.Flags().GetString("profile"))
			slog.Debug("s3", "profile", profileName, "host", config.S3.Server.Host, "port", config.S3.Server.Port)

			var profile *credentials.Profile
			var err error
			if profileName == "" {
				profile, err = credentials.LoadDefaultProfile()
			} else {
				profile, err = credentials.LoadProfile(profileName)
			}

			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}

			client := s3.NewS3Client(
				config.S3.Server.Host,
				config.S3.Server.Port,
				profile.AccessKey,
				profile.SecretKey,
				false,
			)

			handler := handlers.NewS3Handler(client)

			switch {
			case args[0] == "ls":
				handler.ListBuckets()

			case args[0] == "ls" && len(args) >= 2:
				handler.ListObjects(args[1])
			}
		},
	}

	cmd.Flags().String("profile", "", "Profile used for authentication")

	cmd.Flags().String("s3-server-host", config.DefaultS3ServerHost, "Listening host for S3-compatible server")
	cmd.Flags().Uint16("s3-server-port", config.DefaultS3ServerPort, "Listening port for S3-compatible server")

	return cmd
}
