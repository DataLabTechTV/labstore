package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/IllumiKnowLabs/labstore/cli/pkg/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/pkg/handlers"
	"github.com/IllumiKnowLabs/labstore/client/pkg/s3"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/spf13/cobra"
)

func NewS3Cmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "s3",
		Short: "S3 client, designed for learning",

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			credentials.Init()

			profileName := helper.Must(cmd.Flags().GetString("profile"))
			slog.Debug(
				"s3",
				"profile", profileName,
				"host", config.App.Server.S3.Address.Host,
				"port", config.App.Server.S3.Address.Port,
			)

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
				cmd.Context(),
				config.App.Server.S3.Address.Host,
				config.App.Server.S3.Address.Port,
				profile.AccessKey,
				profile.SecretKey,
				false,
				config.App.Server.S3.IO.BufferSize,
			)

			handler := handlers.NewS3Handler(client)

			ctx := context.WithValue(cmd.Context(), handlerKeyCtx, handler)
			cmd.SetContext(ctx)
		},
	}

	cmd.PersistentFlags().String("profile", "", "Profile used for authentication")

	cmd.PersistentFlags().String("s3-address-host", config.DefaultS3AddressHost, "Listening host for S3-compatible server")
	cmd.PersistentFlags().Uint16("s3-address-port", config.DefaultS3AddressPort, "Listening port for S3-compatible server")

	cmd.AddCommand(NewServiceCmd())
	cmd.AddCommand(NewBucketsCmd())
	cmd.AddCommand(NewObjectsCmd())

	return cmd
}
