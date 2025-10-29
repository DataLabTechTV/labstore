package main

import (
	"fmt"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   strings.ToLower(config.Name),
		Short: config.Description,
		Long:  fmt.Sprintln(config.Description, "created by", config.Author),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debug := helper.Must(cmd.Flags().GetBool("debug"))
			logger.Init(logger.WithDebugFlag(debug))
			config.Load()
			iam.Load()
		},
	}

	cmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")

	cmd.AddCommand(NewServeCmd())

	return cmd
}
