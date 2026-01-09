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

const (
	MaxWidth = 80
)

var (
	TitleStyle = lipgloss.NewStyle().
			Width(MaxWidth).
			Align(lipgloss.Center).
			Margin(1, 0).
			Bold(true).
			Foreground(ActivePalette.Accent).
			Background(ActivePalette.SurfaceAlt).
			Render

	HelpStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextMuted).
			Render

	SuccessStyle = lipgloss.NewStyle().
			Width(25).
			Align(lipgloss.Left).
			Bold(true).
			Foreground(ActivePalette.Success).
			Render

	SuccessMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextPrimary).
			Render

	DateStyle = lipgloss.NewStyle().
			Width(25).
			Align(lipgloss.Left).
			Foreground(ActivePalette.TextMuted).
			Render

	SizeStyle = lipgloss.NewStyle().
			Width(20).
			Align(lipgloss.Left).
			Foreground(ActivePalette.AccentMuted).
			Render

	DirStyle = lipgloss.NewStyle().
			Width(60).
			Align(lipgloss.Left).
			Foreground(ActivePalette.Accent).
			Render

	FileStyle = lipgloss.NewStyle().
			Width(60).
			Align(lipgloss.Left).
			Foreground(ActivePalette.TextPrimary).
			Render

	ErrCodeStyle = lipgloss.NewStyle().
			Width(25).
			Foreground(ActivePalette.Error).
			Render

	ErrMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.Accent).
			Render
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
			ErrCodeStyle(s3Error.Code),
			ErrMsgStyle(s3Error.Message),
		)
	case errors.As(err, &iamError):
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle(iamError.Code),
			ErrMsgStyle(iamError.Message),
		)
	default:
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle("Error"),
			ErrMsgStyle(err.Error()),
		)
	}

	fmt.Println(errView)
}

func PrintTitle(title string) {
	titleView := TitleStyle(title)
	fmt.Println(titleView)
}

func PrintSuccess(msg string) {
	successView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		SuccessStyle("✅ Success"),
		SuccessMsgStyle(msg),
	)
	fmt.Println(successView)
}

func PrintBucket(bucket t.Bucket) {
	bucketView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(bucket.CreationDate)),
		DirStyle(fmt.Sprintf("%s/", bucket.Name)),
	)

	fmt.Println(bucketView)
}

func PrintCommonPrefix(commonPrefix t.CommonPrefix) {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(t.Timestamp(time.Now()))),
		SizeStyle(SizeFormat(0)),
		DirStyle(commonPrefix.Prefix),
	)

	fmt.Println(objectView)
}

func PrintObject(object t.Object) {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(object.LastModified)),
		SizeStyle(SizeFormat(object.Size)),
		FileStyle(object.Key),
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
