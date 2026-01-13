package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
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
		helper.CheckFatal(rootCmd.Help())
	}
}
