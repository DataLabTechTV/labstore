package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintError(err error) {
	// TODO: add support for S3 and IAM errors
	// TODO: replace with lipgloss
	fmt.Println(err.Error())
}

func PrintFileList(paths []string) {
	fileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("201")).
		Bold(false)

	for _, path := range paths {
		fmt.Println(fileStyle.Render(path))
	}
}
