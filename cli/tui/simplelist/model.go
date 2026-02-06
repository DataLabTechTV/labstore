package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ParentIndex   int
	InfoPaneIndex *int
	Provider      providers.Provider
	Entries       []string
	Width         int
	Height        int

	table    table.Model
	hCellPad int
}

type SimpleListOption func(m *Model)

func New(parentIndex int, provider providers.Provider, opts ...SimpleListOption) Model {
	tableStyle := table.DefaultStyles()

	tableStyle.Selected = tableStyle.Selected.
		Foreground(render.ActivePalette.Accent).
		Background(render.ActivePalette.SurfaceAlt).
		Bold(false)

	simpleTable := table.New(
		table.WithWidth(0),
		table.WithHeight(0),
		table.WithFocused(false),
		table.WithStyles(tableStyle),
	)

	model := Model{
		ParentIndex: parentIndex,
		Provider:    provider,
		table:       simpleTable,
		hCellPad:    tableStyle.Cell.GetHorizontalPadding(),
	}

	for _, opt := range opts {
		opt(&model)
	}

	return model
}

func WithInfoPaneIndex(infoPaneIndex int) SimpleListOption {
	return func(m *Model) {
		m.InfoPaneIndex = &infoPaneIndex
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return messages.PaneMsg{
			Index: m.ParentIndex,
			Msg:   messages.RefreshMsg{},
		}
	}
}
