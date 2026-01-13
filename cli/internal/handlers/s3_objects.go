package handlers

import (
	"fmt"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/display"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
)

func (h *S3Handler) PutObject(bucket, key, localPath string, debug bool) {
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

	progressBar, err := tui.NewProgressBarModel(debug)
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

func (h *S3Handler) HeadObject(bucket, key string) {
	tui.PrintTitle(fmt.Sprintf("HeadObject: %s/%s", bucket, key))

	out, err := h.Client.HeadObject(bucket, key)
	if err != nil {
		tui.PrintError(err)
		return
	}

	meta := map[string]display.Meta{
		"Content Type":   {Type: display.MetaTypeString, Value: out.ContentType},
		"Content Length": {Type: display.MetaTypeSize, Value: out.ContentLength},
		"Last Modified":  {Type: display.MetaTypeDate, Value: out.LastModified},
	}

	tui.PrintMetadata(out.StatusCode, meta)
}
