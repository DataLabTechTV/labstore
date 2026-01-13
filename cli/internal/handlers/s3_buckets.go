package handlers

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
)

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
