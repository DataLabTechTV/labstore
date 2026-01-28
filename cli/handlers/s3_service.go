package handlers

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/render"
)

func (h *S3Handler) ListBuckets() error {
	fmt.Println(render.Title("ListBuckets"))

	buckets, err := h.Client.ListBuckets()
	if err != nil {
		fmt.Println(render.Error(err))
		return &errs.RuntimeError{}
	}

	for _, bucket := range buckets {
		fmt.Println(render.Bucket(bucket))
	}

	return nil
}
