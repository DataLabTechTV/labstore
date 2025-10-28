package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/server/handlers"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start S3-compatible Lab Store server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Welcome to %s, by %s", config.Name, config.Author)
		handlers.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
