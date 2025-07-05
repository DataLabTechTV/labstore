/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/DataLabTechTV/labstore/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Welcome to Lab Store, by https://youtube.com/@DataLabTechTV")
		log.Debug("Running in debug mode")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
