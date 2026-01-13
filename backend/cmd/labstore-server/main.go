package main

import "github.com/IllumiKnowLabs/labstore/backend/pkg/helper"

func main() {
	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		helper.CheckFatal(rootCmd.Help())
	}
}
