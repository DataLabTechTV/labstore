package render

import (
	"errors"
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"github.com/IllumiKnowLabs/labstore/cli/internal/format"
	"github.com/charmbracelet/lipgloss"
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
			Width(80).
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

func HttpStatus(code int, msg string) string {
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

	return view
}

func Error(err error) string {
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

	return errView
}

func HttpStatusOrError(code int, err error) string {
	if code == 0 {
		return Error(err)
	}
	return HttpStatus(code, err.Error())
}

func Title(title string) string {
	titleView := TitleStyle(title)
	return titleView
}

func Bucket(bucket types.Bucket) string {
	bucketView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(format.Date(bucket.CreationDate)),
		DirStyle(fmt.Sprintf("%s/", bucket.Name)),
	)

	return bucketView
}

func CommonPrefix(commonPrefix types.CommonPrefix) string {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(format.Date(types.Timestamp(time.Now()))),
		SizeStyle(format.Size(0)),
		DirStyle(commonPrefix.Prefix),
	)

	return objectView
}

func Object(object types.Object) string {
	objectView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		DateStyle(format.Date(object.LastModified)),
		SizeStyle(format.Size(object.Size)),
		FileStyle(object.Key),
	)

	return objectView
}

func Metadata(code int, meta map[string]Meta) string {
	return ""
}
