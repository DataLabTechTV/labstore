package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/DataLabTechTV/labstore/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start S3-compatible Lab Store server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Welcome to Lab Store, by https://youtube.com/@DataLabTechTV")
		log.Debug("Running in debug mode")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
