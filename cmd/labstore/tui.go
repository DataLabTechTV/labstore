package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui"
	"github.com/spf13/cobra"
)

func NewTUICmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tui",
		Short: "TUI and helper commands",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			tui.Run(ctx)
		},
	}

	cmd.AddCommand(NewPaletteCmd())

	return cmd
}

func NewPaletteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "palette",
		Short: "Color palette management",
	}

	cmd.AddCommand(NewPalettePreviewCmd())

	return cmd
}

func NewPalettePreviewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "view",
		Short: "View palette",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(render.ActivePalette.Render())
		},
	}

	return cmd
}
