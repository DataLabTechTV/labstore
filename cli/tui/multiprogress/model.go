package multiprogress

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/progress"
)

type Model struct {
	Progress progress.Model

	Current []int64
	Total   int64

	UploadProgressCh   <-chan messages.UploadProgressMsg
	DownloadProgressCh <-chan messages.DownloadProgressMsg
}

func New(numItems int) *Model {
	return &Model{
		Progress:           progress.New(progress.WithDefaultGradient()),
		Current:            make([]int64, numItems),
		UploadProgressCh:   make(chan messages.UploadProgressMsg),
		DownloadProgressCh: make(chan messages.DownloadProgressMsg),
	}
}
