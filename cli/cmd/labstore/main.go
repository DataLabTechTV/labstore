package main

import (
	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func main() {
	rootCmd := NewRootCmd()
	rootCmd.Execute()
}
