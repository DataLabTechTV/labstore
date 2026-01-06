package main

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/logger"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   fmt.Sprintf("%s", strings.ToLower(constants.Name)),
		Short: fmt.Sprintf("%s, by %s", constants.Name, constants.Author),
		Long:  fmt.Sprintf("%s - %s, by %s", constants.Name, constants.Description, constants.Author),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			welcomeMsg := fmt.Sprintf("🚀 Welcome to %s, by %s", constants.Name, constants.Author)
			if helper.SupportsBox() {
				helper.Box(welcomeMsg)
			} else {
				fmt.Println(welcomeMsg)
			}

			debug := helper.Must(cmd.Flags().GetBool("debug"))
			logger.Init(logger.WithDebugFlag(debug))

			config.DisplayDefaultAdminSecretKey = false
			config.Load(cmd)
		},
	}

	cmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")

	cmd.AddCommand(NewS3Cmd())
	cmd.AddCommand(NewIAMCmd())
	cmd.AddCommand(NewAdminCmd())
	cmd.AddCommand(NewTUICmd())

	return cmd
}
