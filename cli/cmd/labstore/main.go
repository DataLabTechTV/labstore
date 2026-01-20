package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/IllumiKnowLabs/labstore/cli/internal/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func main() {
	rootCmd := NewRootCmd()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		var runtimeError *errs.RuntimeError
		if errors.As(err, &runtimeError) {
			os.Exit(2)
		}

		if cmd, _, cmdErr := rootCmd.Find(os.Args[1:]); cmdErr == nil {
			cmd.PrintErrln(err)
			helper.CheckFatal(cmd.Usage())
			os.Exit(1)
		}
	}
}
