package handlers

import (
	"fmt"
	"os"
	"time"

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

	metadata := map[string]render.Meta{
		"Content Type":   render.NewMetaString(out.ContentType),
		"Content Length": render.NewMetaSize(out.ContentLength),
		"Last Modified":  render.NewMetaDate(out.LastModified),
	}

	fmt.Println(render.Metadata(out.StatusCode, metadata))
}

func (h *S3Handler) GetObject(bucket, key, localPath string, debug bool) {
	fmt.Println(render.Title(fmt.Sprintf("GetObject: %s/%s", bucket, key)))

	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}
	defer helper.CloseWithErr(file, &err)

	progressBar, err := tui.NewProgressBarModel(debug)
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}
	go progressBar.Run()
	defer progressBar.Close()

	code, err := h.Client.GetObject(bucket, key, file, progressBar.Progress)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
	}
	time.Sleep(50 * time.Millisecond)

	fmt.Println(render.HttpStatus(code, fmt.Sprintf("%s downloaded", localPath)))
}

func (h *S3Handler) DeleteObjects(bucket string, keys ...string) {
	if len(keys) == 1 {
		fmt.Println(render.Title(fmt.Sprintf("DeleteObject: %s/%s", bucket, keys[0])))
	} else {
		fmt.Println(render.Title(fmt.Sprintf("DeleteObjects: %s (%d objects)", bucket, len(keys))))
	}

	res, code, err := h.Client.DeleteObjects(bucket, keys...)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
	}

	okCount := 0
	for _, deleted := range res.Deleted {
		fmt.Println(render.HttpStatus(code, deleted.Key))
		okCount++
	}

	for _, s3Error := range res.Error {
		s3Error.StatusCode = code
		fmt.Println(render.Error(&s3Error))
	}

	fmt.Println(render.HttpStatus(code, fmt.Sprintf("%d out of %d object(s) deleted", okCount, len(keys))))
}
