package render

import (
	"slices"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Number interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32
}

type Metadata map[string]Meta

type Meta struct {
	Value  any
	Format func() string
}

func NewString(value string) Meta {
	return Meta{
		Value: value,
		Format: func() string {
			return value
		},
	}
}

func NewDate(value time.Time) Meta {
	return Meta{
		Value: value,
		Format: func() string {
			return time.Time(value).Format(types.ISO8601)
		},
	}
}

func NewSize(value int64) Meta {
	return Meta{
		Value: value,
		Format: func() string {
			p := message.NewPrinter(language.English)
			return p.Sprintf("%d B", value)
		},
	}
}

func NewNumber[T Number](value T) Meta {
	return Meta{
		Value: value,
		Format: func() string {
			p := message.NewPrinter(language.English)

			switch val := any(value).(type) {
			case float32, float64:
				return p.Sprintf("%.6f", val)
			default:
				return p.Sprintf("%d", val)
			}
		},
	}
}

func (metadata Metadata) Render() string {
	metaLabelStyle := lipgloss.NewStyle().
		Width(20).
		Bold(true).
		Align(lipgloss.Right).
		PaddingRight(1).
		MarginRight(2).
		Background(ActivePalette.Surface).
		Foreground(ActivePalette.TextPrimary).
		Render

	MetaValueStyle := lipgloss.NewStyle().
		Foreground(ActivePalette.AccentMuted).
		Render

	rows := []string{}

	labels := make([]string, 0, len(metadata))
	for label := range metadata {
		labels = append(labels, label)
	}
	slices.Sort(labels)

	for _, label := range labels {
		metaRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			metaLabelStyle(label),
			MetaValueStyle(metadata[label].Format()),
		)
		rows = append(rows, metaRow)
	}

	metaView := lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)

	return metaView + "\n"
}
