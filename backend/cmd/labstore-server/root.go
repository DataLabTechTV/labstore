package main

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/constants"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
	"github.com/IllumiKnowLabs/labstore/backend/internal/logger"
	"github.com/IllumiKnowLabs/labstore/backend/internal/profiler"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   fmt.Sprintf("%s-server", strings.ToLower(constants.Name)),
		Short: fmt.Sprintf("%s, by %s", constants.Name, constants.Author),
		Long:  fmt.Sprintf("%s - %s, by %s", constants.Name, constants.Description, constants.Author),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			helper.Box(fmt.Sprintf("🚀 Welcome to %s, by %s", constants.Name, constants.Author))

			debug := helper.Must(cmd.Flags().GetBool("debug"))
			logger.Init(logger.WithDebugFlag(debug))

			config.Load(cmd)
			iam.Load()

			if run_pprof := helper.Must(cmd.Flags().GetBool("pprof")); run_pprof {
				pprof_host := helper.Must(cmd.Flags().GetString("pprof-host"))
				pprof_port := helper.Must(cmd.Flags().GetInt("pprof-port"))

				pprof := profiler.NewProfiler(
					profiler.WithHost(pprof_host),
					profiler.WithPort(pprof_port),
				)

				pprof.Start()
			}
		},
	}

	cmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")
	cmd.PersistentFlags().Bool("pprof", false, "Enable profiler")
	cmd.PersistentFlags().String("pprof-host", "localhost", "Profiler host")
	cmd.PersistentFlags().Int("pprof-port", 6060, "Profiler port")

	cmd.AddCommand(NewServeCmd())

	return cmd
}
