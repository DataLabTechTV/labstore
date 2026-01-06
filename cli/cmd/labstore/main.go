package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func main() {
	rootCmd := NewRootCmd()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rootCmd.ExecuteContext(ctx)
}
