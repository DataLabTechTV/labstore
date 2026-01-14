package handlers

import (
	"fmt"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
)

func (h *S3Handler) PutObject(bucket, key, localPath string, debug bool) {
	if !helper.FileExists(localPath) {
		fmt.Println(render.Error(fmt.Errorf("file not found: %s", localPath)))
		return
	}

	fmt.Println(render.Title(fmt.Sprintf("PutObject: %s/%s", bucket, key)))

	file, err := os.Open(localPath)
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}

	progressBar, err := tui.NewProgressBarModel(debug)
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}
	go progressBar.Run()
	defer progressBar.Close()

	code, err := h.Client.PutObject(bucket, key, file, progressBar.Progress)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
	}

	fmt.Println(render.HttpStatus(code, fmt.Sprintf("%s uploaded", localPath)))
}

func (h *S3Handler) HeadObject(bucket, key string) {
	fmt.Println(render.Title(fmt.Sprintf("HeadObject: %s/%s", bucket, key)))

	out, err := h.Client.HeadObject(bucket, key)
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}

	meta := map[string]render.Meta{
		"Content Type":   {Type: render.MetaTypeString, Value: out.ContentType},
		"Content Length": {Type: render.MetaTypeSize, Value: out.ContentLength},
		"Last Modified":  {Type: render.MetaTypeDate, Value: out.LastModified},
	}

	fmt.Println(render.Metadata(out.StatusCode, meta))
}
