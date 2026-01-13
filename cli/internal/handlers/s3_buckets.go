package handlers

import (
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
)

func (h *S3Handler) CreateBucket(bucket string) {
	tui.PrintTitle(fmt.Sprintf("CreateBucket: %s", bucket))

	code, err := h.Client.CreateBucket(bucket)
	if err != nil {
		tui.PrintStatusOrError(code, err)
		return
	}

	tui.PrintStatus(code, "Bucket created")
}

func (h *S3Handler) HeadBucket(bucket string) {
	tui.PrintTitle(fmt.Sprintf("HeadBucket: %s", bucket))

	code, err := h.Client.HeadBucket(bucket)
	if err != nil {
		tui.PrintStatusOrError(code, err)
		return
	}

	switch code {
	case http.StatusForbidden:
		tui.PrintStatus(code, "Bucket access denied")
	case http.StatusNotFound:
		tui.PrintStatus(code, "Bucket does not exist")
	case http.StatusOK:
		tui.PrintStatus(code, "Bucket access allowed")
	default:
		tui.PrintStatus(code, "Unknown status")
	}
}

func (h *S3Handler) ListObjects(bucket string, key *string) {
	if key == nil {
		key = helper.StringPtr("/")
	}

	tui.PrintTitle(fmt.Sprintf("ListObjectsV2: %s", bucket))

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			tui.PrintError(r.Err)
			return
		}

		switch {
		case r.IsObject():
			tui.PrintObject(*r.Object)
		case r.IsCommonPrefix():
			tui.PrintCommonPrefix(*r.CommonPrefix)
		}
	}
}
