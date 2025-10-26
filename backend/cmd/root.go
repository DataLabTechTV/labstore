package cmd

import (
	"fmt"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/config"
	"github.com/DataLabTechTV/labstore/backend/iam"
	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   strings.ToLower(config.Name),
	Short: config.Description,
	Long:  fmt.Sprintln(config.Description, "created by", config.Author),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug := helper.Must(cmd.Flags().GetBool("debug"))

		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})

		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug mode: on")
		} else {
			log.SetLevel(log.InfoLevel)
		}

		config.LoadEnv()
		iam.Load()
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")
}

func Execute() {
	helper.Check(rootCmd.Execute())
}
