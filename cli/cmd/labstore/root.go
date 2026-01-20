package main

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/iam"
	"github.com/IllumiKnowLabs/labstore/server/pkg/logger"
	"github.com/IllumiKnowLabs/labstore/server/pkg/profiler"
	"github.com/spf13/cobra"
)

type contextKey string

const handlerKeyCtx contextKey = "handler"

func NewRootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   strings.ToLower(constants.Name),
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

			iam.Init()

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

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")
	cmd.PersistentFlags().Bool("pprof", false, "Enable profiler")
	cmd.PersistentFlags().String("pprof-host", "localhost", "Profiler host")
	cmd.PersistentFlags().Int("pprof-port", 6060, "Profiler port")

	cmd.AddCommand(NewServeCmd())
	cmd.AddCommand(NewS3Cmd())
	cmd.AddCommand(NewIAMCmd())
	cmd.AddCommand(NewAdminCmd())
	cmd.AddCommand(NewTUICmd())

	return cmd
}
