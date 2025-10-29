package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/service"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start S3-compatible Lab Store server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Welcome to %s, by %s", config.Name, config.Author)
			service.Start()
		},
	}

	return cmd
}
