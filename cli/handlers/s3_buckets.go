package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/server/helper"
)

func (h *S3Handler) CreateBucket(bucket string) error {
	fmt.Println(render.Title(fmt.Sprintf("CreateBucket: %s", bucket)))

	code, err := h.Client.CreateBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(code, "Bucket created"))
	return nil
}

func (h *S3Handler) HeadBucket(bucket string) error {
	fmt.Println(render.Title(fmt.Sprintf("HeadBucket: %s", bucket)))

	code, err := h.Client.HeadBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
	}

	switch code {
	case http.StatusForbidden:
		fmt.Println(render.HttpStatus(code, "Bucket access denied"))
	case http.StatusNotFound:
		fmt.Println(render.HttpStatus(code, "Bucket does not exist"))
	case http.StatusOK:
		fmt.Println(render.HttpStatus(code, "Bucket access allowed"))
	default:
		fmt.Println(render.HttpStatus(code, "Unknown status"))
	}

	return nil
}

func (h *S3Handler) ListObjects(bucket string, key *string) error {
	if key == nil {
		key = helper.StringPtr("/")
	}

	fmt.Println(render.Title(fmt.Sprintf("ListObjectsV2: %s", bucket)))

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			fmt.Println(render.Error(r.Err))
		}

		switch {
		case r.IsObject():
			fmt.Println(render.Object(*r.Object))
		case r.IsCommonPrefix():
			fmt.Println(render.CommonPrefix(*r.CommonPrefix))
		}
	}

	return nil
}

func (h *S3Handler) DeleteBucket(bucket string) error {
	fmt.Println(render.Title(fmt.Sprintf("DeleteBucket: %s", bucket)))

	code, err := h.Client.DeleteBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return &errs.RuntimeError{}
	}

	fmt.Println(render.HttpStatus(code, "Bucket deleted"))
	return nil
}
