package tui

import (
	"errors"
	"fmt"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
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

	ErrCodeStyle = lipgloss.NewStyle().
			Width(25).
			Foreground(ActivePalette.Error)

	ErrMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.Accent)
)

func PrintError(err error) {
	var (
		s3Error  *errs.S3Error
		iamError *errs.IAMError
		errView  string
	)

	switch {
	case errors.As(err, &s3Error):
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle.Render(s3Error.Code),
			ErrMsgStyle.Render(s3Error.Message),
		)
	case errors.As(err, &iamError):
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle.Render(iamError.Code),
			ErrMsgStyle.Render(iamError.Message),
		)
	default:
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle.Render("Error"),
			ErrMsgStyle.Render(err.Error()),
		)
	}

	fmt.Println(errView)
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

func PrintObject(bucket string, object t.Object) {
	path := fmt.Sprintf("%s/%s", bucket, object.Key)
	objectView := PathStyle.Render(path)
	fmt.Println(objectView)
}
