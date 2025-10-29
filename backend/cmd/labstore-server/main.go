package main

import "github.com/DataLabTechTV/labstore/backend/internal/helper"

func main() {
	rootCmd := NewRootCmd()
	helper.CheckFatal(rootCmd.Execute())
}
