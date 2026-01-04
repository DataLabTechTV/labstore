package tui

import (
	"fmt"

	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/charmbracelet/lipgloss"
)

var (
	DateStyle = lipgloss.NewStyle().
			Width(25).
			Align(lipgloss.Left).
			Foreground(ActivePalette.TextMuted)

	PathStyle = lipgloss.NewStyle().
			Width(60).
			Align(lipgloss.Left).
			Foreground(ActivePalette.TextPrimary)
)

func PrintError(err error) {
	// TODO: add support for S3 and IAM errors
	// TODO: replace with lipgloss
	fmt.Println(err.Error())
}

func PrintBuckets(buckets []t.Bucket) {

	for _, bucket := range buckets {
		date := fmt.Sprintf("[%s]", bucket.CreationDate)
		bucketView := lipgloss.JoinHorizontal(
			lipgloss.Top,
			DateStyle.Render(date),
			PathStyle.Render(bucket.Name+"/"),
		)

		fmt.Println(bucketView)
	}
}

func PrintObjects(bucket string, objects []t.Object) {
	for _, object := range objects {
		path := fmt.Sprintf("%s/%s", bucket, object.Key)
		objectView := PathStyle.Render(path)
		fmt.Println(objectView)
	}
}
