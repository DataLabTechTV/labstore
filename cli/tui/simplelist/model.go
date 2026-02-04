package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Provider providers.Provider
	Entries  []string
	Width    int
	Height   int

	table    table.Model
	hCellPad int
}

func New(provider providers.Provider) Model {
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

	return Model{
		Provider: provider,
		table:    simpleTable,
		hCellPad: tableStyle.Cell.GetHorizontalPadding(),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return messages.SimpleListMsg{Msg: messages.RefreshMsg{}}
	}
}
