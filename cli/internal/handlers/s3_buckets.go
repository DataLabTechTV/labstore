package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
)

func (h *S3Handler) CreateBucket(bucket string) {
	fmt.Println(render.Title(fmt.Sprintf("CreateBucket: %s", bucket)))

	code, err := h.Client.CreateBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
	}

	fmt.Println(render.HttpStatus(code, "Bucket created"))
}

func (h *S3Handler) HeadBucket(bucket string) {
	fmt.Println(render.Title(fmt.Sprintf("HeadBucket: %s", bucket)))

	code, err := h.Client.HeadBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
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
}

func (h *S3Handler) ListObjects(bucket string, key *string) {
	if key == nil {
		key = helper.StringPtr("/")
	}

	fmt.Println(render.Title(fmt.Sprintf("ListObjectsV2: %s", bucket)))

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			fmt.Println(render.Error(r.Err))
			return
		}

		switch {
		case r.IsObject():
			fmt.Println(render.Object(*r.Object))
		case r.IsCommonPrefix():
			fmt.Println(render.CommonPrefix(*r.CommonPrefix))
		}
	}
}

func (h *S3Handler) DeleteBucket(bucket string) {
	fmt.Println(render.Title(fmt.Sprintf("DeleteBucket: %s", bucket)))

	code, err := h.Client.DeleteBucket(bucket)
	if err != nil {
		fmt.Println(render.HttpStatusOrError(code, err))
		return
	}

	fmt.Println(render.HttpStatus(code, "Bucket deleted"))
}
