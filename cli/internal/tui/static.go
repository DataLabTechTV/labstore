package tui

import (
	"errors"
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/IllumiKnowLabs/labstore/cli/internal/display"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	MaxWidth = 80
)

var (
	codeStyle = lipgloss.NewStyle().
			MarginRight(2).
			Bold(true)

	InfoCodeStyle = codeStyle.
			Foreground(ActivePalette.TextPrimary).
			Render

	InfoMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextPrimary).
			Render

	SuccessCodeStyle = codeStyle.
				Foreground(ActivePalette.Success).
				Render

	SuccessMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextPrimary).
			Render

	WarnCodeStyle = codeStyle.
			Foreground(ActivePalette.Warning).
			Render

	WarnMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.TextPrimary).
			Render

	ErrCodeStyle = codeStyle.
			Foreground(ActivePalette.Error).
			Render

	ErrMsgStyle = lipgloss.NewStyle().
			Foreground(ActivePalette.Accent).
			Render

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

	DateStyle = lipgloss.NewStyle().
			Width(25).
			Foreground(ActivePalette.TextMuted).
			Render

	SizeStyle = lipgloss.NewStyle().
			Width(20).
			Foreground(ActivePalette.AccentMuted).
			Render

	DirStyle = lipgloss.NewStyle().
			Width(60).
			Foreground(ActivePalette.Accent).
			Render

	FileStyle = lipgloss.NewStyle().
			Width(60).
			Foreground(ActivePalette.TextPrimary).
			Render
)

func PrintStatus(code int, msg string) {
	var view string

	switch {
	case code >= 200 && code < 300:
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			SuccessCodeStyle(fmt.Sprintf("[%d OK]", code)),
			SuccessMsgStyle(msg),
		)
	case code >= 300 && code < 400:
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			WarnCodeStyle(fmt.Sprintf("[%d Redirect]", code)),
			WarnMsgStyle(msg),
		)
	case code >= 400 && code < 600:
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle(fmt.Sprintf("[%d Error]", code)),
			ErrMsgStyle(msg),
		)
	default:
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			InfoCodeStyle(fmt.Sprintf("[%d Info]", code)),
			InfoMsgStyle(msg),
		)
	}

	fmt.Println(view)
}

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
			ErrCodeStyle(fmt.Sprintf("[%d S3 Error] %s", s3Error.StatusCode, s3Error.Code)),
			ErrMsgStyle(s3Error.Message),
		)
	case errors.As(err, &iamError):
		errView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ErrCodeStyle(fmt.Sprintf("[%d IAM Error] %s", iamError.StatusCode, iamError.Code)),
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

func PrintStatusOrError(code int, err error) {
	if code == 0 {
		PrintError(err)
	} else {
		PrintStatus(code, err.Error())
	}
}

func PrintTitle(title string) {
	titleView := TitleStyle(title)
	fmt.Println(titleView)
}

func PrintBucket(bucket types.Bucket) {
	bucketView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(bucket.CreationDate)),
		DirStyle(fmt.Sprintf("%s/", bucket.Name)),
	)

	fmt.Println(bucketView)
}

func PrintCommonPrefix(commonPrefix types.CommonPrefix) {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(types.Timestamp(time.Now()))),
		SizeStyle(SizeFormat(0)),
		DirStyle(commonPrefix.Prefix),
	)

	fmt.Println(objectView)
}

func PrintObject(object types.Object) {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(DateFormat(object.LastModified)),
		SizeStyle(SizeFormat(object.Size)),
		FileStyle(object.Key),
	)

	fmt.Println(objectView)
}

func PrintMetadata(code int, meta map[string]display.Meta) {

}

func SizeFormat(size int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d B", size)
}

func DateFormat(date types.Timestamp) string {
	return fmt.Sprintf("[%s]", time.Time(date).Format(types.ISO8601))
}
