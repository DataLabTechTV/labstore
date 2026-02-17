package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/ui/progressbar"
	"github.com/IllumiKnowLabs/labstore/server/helper"
)

func (h *S3Handler) PutObject(bucket, key, localPath string, debug bool) error {
	if !helper.FileExists(localPath) {
		err := fmt.Errorf("file not found: %s", localPath)
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.Title(fmt.Sprintf("PutObject: %s/%s", bucket, key)))

	file, err := os.Open(localPath)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	progressBar, err := progressbar.New(h.Client.Ctx, debug)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}
	go progressBar.Run()
	defer progressBar.Close()

	code, err := h.Client.PutObject(bucket, key, file, progressBar.Progress)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(code, fmt.Sprintf("%s uploaded", localPath)))
	return nil
}

func (h *S3Handler) HeadObject(bucket, key string) error {
	fmt.Println(render.Title(fmt.Sprintf("HeadObject: %s/%s", bucket, key)))

	out, err := h.Client.HeadObject(bucket, key)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	metadata := render.Metadata{
		"Content Type":   render.NewString(out.ContentType),
		"Content Length": render.NewSize(out.ContentLength),
		"Last Modified":  render.NewDate(out.LastModified),
	}

	fmt.Println(metadata.Render())
	return nil
}

func (h *S3Handler) GetObject(bucket, key, localPath string, debug bool) error {
	fmt.Println(render.Title(fmt.Sprintf("GetObject: %s/%s", bucket, key)))

	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}
	defer helper.CloseWithErr(file, &err)

	progressBar, err := progressbar.New(h.Client.Ctx, debug)
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}
	go progressBar.Run()
	defer progressBar.Close()

	code, err := h.Client.GetObject(bucket, key, file, progressBar.Progress)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
	}
	time.Sleep(50 * time.Millisecond)

	fmt.Println(render.HttpStatus(code, fmt.Sprintf("%s downloaded", localPath)))
	return nil
}

func (h *S3Handler) DeleteObjects(bucket string, keys ...string) error {
	if len(keys) == 1 {
		fmt.Println(render.Title(fmt.Sprintf("DeleteObject: %s/%s", bucket, keys[0])))
	} else {
		fmt.Println(render.Title(fmt.Sprintf("DeleteObjects: %s (%d objects)", bucket, len(keys))))
	}

	res, code, err := h.Client.DeleteObjects(bucket, keys...)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
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
	return nil
}
