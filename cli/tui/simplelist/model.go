package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
)

type Model struct {
	Entries []string
	Width   int
	Height  int

	table    table.Model
	hCellPad int
}

func New() Model {
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
		table:    simpleTable,
		hCellPad: tableStyle.Cell.GetHorizontalPadding(),
	}

	return model
}

func (m Model) SetEntries(entries []providers.Entry) Model {
	m.Entries = providers.EntryNames(entries)
	m.updateTable()
	return m
}

func (m Model) Clear() Model {
	m.Entries = []string{}
	m.table.SetRows([]table.Row{})
	m.table.GotoTop()
	return m
}

func (m Model) Selected() (string, bool) {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) < 1 {
		return "", false
	}
	value := selectedRow[0]
	return value, true
}
