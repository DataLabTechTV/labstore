package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
)

func (h *S3Handler) CreateBucket(bucket string) {
	render.Title(fmt.Sprintf("CreateBucket: %s", bucket))

	code, err := h.Client.CreateBucket(bucket)
	if err != nil {
		render.HttpStatusOrError(code, err)
		return
	}

	render.HttpStatus(code, "Bucket created")
}

func (h *S3Handler) HeadBucket(bucket string) {
	render.Title(fmt.Sprintf("HeadBucket: %s", bucket))

	code, err := h.Client.HeadBucket(bucket)
	if err != nil {
		render.HttpStatusOrError(code, err)
		return
	}

	switch code {
	case http.StatusForbidden:
		render.HttpStatus(code, "Bucket access denied")
	case http.StatusNotFound:
		render.HttpStatus(code, "Bucket does not exist")
	case http.StatusOK:
		render.HttpStatus(code, "Bucket access allowed")
	default:
		render.HttpStatus(code, "Unknown status")
	}
}

func (h *S3Handler) ListObjects(bucket string, key *string) {
	if key == nil {
		key = helper.StringPtr("/")
	}

	render.Title(fmt.Sprintf("ListObjectsV2: %s", bucket))

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			render.Error(r.Err)
			return
		}

		switch {
		case r.IsObject():
			render.Object(*r.Object)
		case r.IsCommonPrefix():
			render.CommonPrefix(*r.CommonPrefix)
		}
	}
}

func (h *S3Handler) DeleteBucket(bucket string) {
	render.Title(fmt.Sprintf("DeleteBucket: %s", bucket))

	code, err := h.Client.DeleteBucket(bucket)
	if err != nil {
		render.HttpStatusOrError(code, err)
		return
	}

	render.HttpStatus(code, "Bucket deleted")
}
