package handlers

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/internal/render"
)

func (h *S3Handler) ListBuckets() {
	fmt.Println(render.Title("ListBuckets"))

	buckets, err := h.Client.ListBuckets()
	if err != nil {
		fmt.Println(render.Error(err))
		return
	}

	for _, bucket := range buckets {
		fmt.Println(render.Bucket(bucket))
	}
}
