package tui

import "fmt"

func PrintError(err error) {
	// TODO: add support for S3 and IAM errors
	// TODO: replace with lipgloss
	fmt.Println(err.Error())
}

func PrintFileList(paths []string) {
	// TODO: replace with lipgloss
	for _, path := range paths {
		fmt.Println(path)
	}
}
