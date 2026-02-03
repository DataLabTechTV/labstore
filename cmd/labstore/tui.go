package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui"
	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/spf13/cobra"
)

func NewTUICmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tui",
		Short: "TUI and helper commands",
		Annotations: map[string]string{
			"skip-bootstrap": "yes",
		},
		Run: func(cmd *cobra.Command, args []string) {
			config.Load(cmd)

			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			tui.Run(ctx)
		},
	}

	cmd.Flags().String("iam-address-host", config.DefaultIAMAddressHost, "Host for IAM endpoint")
	cmd.Flags().String("iam-address-port", config.DefaultIAMAddressHost, "Port for IAM endpoint")

	cmd.Flags().String("s3-address-host", config.DefaultS3AddressHost, "Host for S3-compatible endpoint")
	cmd.Flags().Uint16("s3-address-port", config.DefaultS3AddressPort, "Port for S3-compatible endpoint")

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
