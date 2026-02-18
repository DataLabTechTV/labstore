package multiprogress

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/progress"
)

type Model struct {
	Progress progress.Model
	Width    int
	Height   int

	Current []int64
	Total   int64

	UploadProgressCh   <-chan messages.UploadProgressMsg
	DownloadProgressCh <-chan messages.DownloadProgressMsg
}

func New(numItems, width, height int) *Model {
	return &Model{
		Progress:           progress.New(progress.WithDefaultGradient()),
		Width:              width,
		Height:             height,
		Current:            make([]int64, numItems),
		UploadProgressCh:   make(chan messages.UploadProgressMsg),
		DownloadProgressCh: make(chan messages.DownloadProgressMsg),
	}
}
