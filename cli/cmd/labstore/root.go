package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		// Use:   fmt.Sprintf("%s", strings.ToLower(constants.Name)),
		// Short: fmt.Sprintf("%s, by %s", constants.Name, constants.Author),
		// Long:  fmt.Sprintf("%s - %s, by %s", constants.Name, constants.Description, constants.Author),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("TODO")
		},
	}

	cmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")
	cmd.PersistentFlags().Bool("pprof", false, "Enable profiler")
	cmd.PersistentFlags().String("pprof-host", "localhost", "Profiler host")
	cmd.PersistentFlags().Int("pprof-port", 6060, "Profiler port")

	return cmd
}
