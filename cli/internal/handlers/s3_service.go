package handlers

import "github.com/IllumiKnowLabs/labstore/cli/internal/render"

func (h *S3Handler) ListBuckets() {
	render.Title("ListBuckets")

	buckets, err := h.Client.ListBuckets()
	if err != nil {
		render.Error(err)
		return
	}

	for _, bucket := range buckets {
		render.Bucket(bucket)
	}
}
