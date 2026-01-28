package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/pkg/render"
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
			if baseCmd := topLevelCommand(cmd); baseCmd.Name() == "completion" {
				return
			}

			if cmd.Annotations["show-default-secret"] == "yes" {
				config.DisplayDefaultAdminSecretKey = true
			}

			if cmd.Annotations["mode"] == "daemon" {
				slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
				config.Load(cmd)
				return
			}

			bootstrap(cmd)
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
	AddDaemonCommands(cmd)

	return cmd
}

func topLevelCommand(cmd *cobra.Command) *cobra.Command {
	for cmd.Parent().HasParent() {
		cmd = cmd.Parent()
	}
	return cmd
}

func bootstrap(cmd *cobra.Command) {
	welcomeMsg := fmt.Sprintf("🚀 Welcome to %s, by %s", constants.Name, constants.Author)
	if render.SupportsBox() {
		fmt.Fprintln(os.Stderr, render.Box(welcomeMsg))
	} else {
		fmt.Fprintln(os.Stderr, welcomeMsg)
	}

	debug := helper.Must(cmd.Flags().GetBool("debug"))
	logger.Init(logger.WithDebugFlag(debug))

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
}
