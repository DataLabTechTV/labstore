package tui

import (
	"errors"
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Width(79).
			Align(lipgloss.Center).
			Margin(1, 0).
			Bold(true).
			Foreground(ActivePalette.Accent).
			Background(ActivePalette.SurfaceAlt)

	SuccessStyle = lipgloss.NewStyle().
			Width(25).
			Align(lipgloss.Left).
			Bold(true).
			Foreground(ActivePalette.Success)

	SuccessMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextPrimary)

	DateStyle = lipgloss.NewStyle().
			Width(25).
			Align(lipgloss.Left).
			Foreground(ActivePalette.TextMuted)

	SizeStyle = lipgloss.NewStyle().
			Width(20).
			Align(lipgloss.Left).
			Foreground(ActivePalette.Accent)

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

func PrintTitle(title string) {
	titleView := TitleStyle.Render(title)
	fmt.Println(titleView)
}

func PrintSuccess(msg string) {
	successView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		SuccessStyle.Render("✅ Success"),
		SuccessMsgStyle.Render(msg),
	)
	fmt.Println(successView)
}

func PrintBucket(bucket t.Bucket) {
	bucketView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle.Render(DateFormat(bucket.CreationDate)),
		PathStyle.Render(bucket.Name+"/"),
	)

	fmt.Println(bucketView)
}

func PrintObject(bucket string, object t.Object) {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle.Render(DateFormat(object.LastModified)),
		SizeStyle.Render(SizeFormat(object.Size)),
		PathStyle.Render(object.Key),
	)

	fmt.Println(objectView)
}

func SizeFormat(size int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d B", size)
}

func DateFormat(date t.Timestamp) string {
	return fmt.Sprintf("[%s]", time.Time(date).Format(t.ISO8601))
}
