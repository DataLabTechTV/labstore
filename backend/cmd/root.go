package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "labstore",
	Short: "Lab Store is a minimal S3-compatible object store",
	Long: heredoc.Doc(`
		Lab Store is a minimal S3-compatible object store,
		created by https://youtube.com/@DataLabTechTV
	`),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var debug, _ = cmd.Flags().GetBool("debug")

		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})

		if debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Set debug level for logging")
}
