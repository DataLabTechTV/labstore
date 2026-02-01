package main

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/server/constants"
	"github.com/IllumiKnowLabs/labstore/server/router"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Display version and build details",
		Annotations: map[string]string{
			"skip-bootstrap": "yes",
		},
		Run: func(cmd *cobra.Command, args []string) {
			printWelcome()

			labelStyle := lipgloss.NewStyle().
				Foreground(render.ActivePalette.Accent).
				Render

			fmt.Printf(
				"%s: %s (%s: %s, %s: %s)\n",
				labelStyle("version"), constants.Version,
				labelStyle("tag"), constants.GitTag,
				labelStyle("commit"), constants.GitCommit,
			)

			fmt.Printf(
				"%s: %s, %s: %s\n",
				labelStyle("build time"), constants.BuildTime,
				labelStyle("builder"), constants.Builder,
			)

			boolTrueStyle := lipgloss.NewStyle().
				Foreground(render.ActivePalette.Success).
				Render

			boolFalseStyle := lipgloss.NewStyle().
				Foreground(render.ActivePalette.Error).
				Render

			var embedAssetsView string
			if router.EmbedAssets {
				embedAssetsView = boolTrueStyle("yes")
			} else {
				embedAssetsView = boolFalseStyle("no")
			}
			fmt.Printf("%s: %v\n", labelStyle("embedded web assets"), embedAssetsView)
		},
	}

	return cmd
}
