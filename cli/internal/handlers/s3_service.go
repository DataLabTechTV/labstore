package handlers

import "github.com/IllumiKnowLabs/labstore/cli/internal/tui"

func (h *S3Handler) ListBuckets() {
	tui.PrintTitle("ListBuckets")

	buckets, err := h.Client.ListBuckets()
	if err != nil {
		tui.PrintError(err)
		return
	}

	for _, bucket := range buckets {
		tui.PrintBucket(bucket)
	}
}
