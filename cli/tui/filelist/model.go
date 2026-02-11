package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
)

type Model struct {
	Entries []providers.Entry
	Width   int
	Height  int

	table    table.Model
	hCellPad int
}

func New() Model {
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
		table:    fileTable,
		hCellPad: tableStyle.Cell.GetHorizontalPadding(),
	}
}

func (m Model) SetEntries(entries []providers.Entry, active *string) Model {
	m.Entries = entries
	m.updateTable(active)
	return m
}

func (m Model) Clear() Model {
	m.Entries = []providers.Entry{}
	m.table.SetRows([]table.Row{})
	m.table.SetCursor(0)
	return m
}

func (m Model) Selected() (string, bool) {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) < 3 {
		return "", false
	}
	value := selectedRow[2]
	return value, true
}
