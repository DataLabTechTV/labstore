package main

import "github.com/IllumiKnowLabs/labstore/backend/pkg/helper"

func main() {
	rootCmd := NewRootCmd()
	helper.CheckFatal(rootCmd.Execute())
}
