package main

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func main() {
	rootCmd := NewRootCmd()
	helper.CheckFatal(rootCmd.Execute())
}
