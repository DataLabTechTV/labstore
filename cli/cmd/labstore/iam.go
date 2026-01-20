package main

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/internal/handlers"
	"github.com/IllumiKnowLabs/labstore/client/iam"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/spf13/cobra"
)

func NewIAMCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "iam",
		Short: "IAM client, designed for learning",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client := iam.NewIAMClient(
				cmd.Context(),
				config.IAM.Address.Host,
				config.IAM.Address.Port,
			)

			handler := handlers.NewIAMHandler(client)

			ctx := context.WithValue(cmd.Context(), handlerKeyCtx, handler)
			cmd.SetContext(ctx)
		},
	}

	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	cmd.Flags().String("iam-server-host", config.DefaultIAMServerHost, "Listening host for IAM server")
	cmd.Flags().String("iam-server-port", config.DefaultIAMServerHost, "Listening port for IAM server")

	cmd.AddCommand(NewUsersCmd())
	cmd.AddCommand(NewGroupsCmd())
	cmd.AddCommand(NewPoliciesCmd())

	return cmd
}
