package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ParentID int
	Provider providers.Provider
	Entries  []providers.Entry
	Width    int
	Height   int

	table    table.Model
	hCellPad int
}

func New(parentID int, provider providers.Provider) Model {
	tableStyle := table.DefaultStyles()

	tableStyle.Header = tableStyle.Header.
		Foreground(render.ActivePalette.SurfaceAlt).
		Bold(true)

	tableStyle.Selected = tableStyle.Selected.
		Foreground(render.ActivePalette.Accent).
		Background(render.ActivePalette.SurfaceAlt).
		Bold(false)

	fileTable := table.New(
		table.WithWidth(0),
		table.WithHeight(0),
		table.WithFocused(false),
		table.WithStyles(tableStyle),
	)

	return Model{
		ParentID: parentID,
		Provider: provider,
		table:    fileTable,
		hCellPad: tableStyle.Cell.GetHorizontalPadding(),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return messages.FileListMsg{Msg: messages.RefreshMsg{}}
	}
}
