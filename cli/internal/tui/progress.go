package tui

import (
	"fmt"
	"os"

	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	tea.Model
	progress progress.Model
}

func NewProgressBar() s3.ProgressCallback {
	// TODO: implement
	m := UI{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	return nil
}
