package handlers

import (
	"fmt"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
)

func (h *S3Handler) PutObject(bucket, key, localPath string) {
	if !helper.FileExists(localPath) {
		tui.PrintError(fmt.Errorf("file not found: %s", localPath))
		return
	}

	tui.PrintTitle(fmt.Sprintf("PutObject: %s/%s", bucket, key))

	file, err := os.Open(localPath)
	if err != nil {
		tui.PrintError(err)
		return
	}

	progressBar, err := tui.NewProgressBarModel()
	if err != nil {
		tui.PrintError(err)
		return
	}
	go progressBar.Run()
	defer progressBar.Close()

	code, err := h.Client.PutObject(bucket, key, file, progressBar.Progress)
	if err != nil {
		tui.PrintStatusOrError(code, err)
		return
	}

	tui.PrintStatus(code, fmt.Sprintf("%s uploaded", localPath))
}
