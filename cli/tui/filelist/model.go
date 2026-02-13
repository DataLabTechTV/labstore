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

func (m *Model) SetEntries(entries []providers.Entry, active *string) {
	m.Entries = entries
	m.updateTable(active)
}

func (m *Model) Clear() {
	m.Entries = []providers.Entry{}
	m.table.SetRows([]table.Row{})
	m.table.GotoTop()
}

func (m Model) Selected() (string, bool) {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) < 4 {
		return "", false
	}
	value := selectedRow[3]
	return value, true
}

func (m *Model) Mark() {
	if idx := m.table.Cursor(); idx < len(m.Entries) {
		m.Entries[idx].Marked = !m.Entries[idx].Marked
		m.updateTable(nil)
	}
}

func (m Model) Marked() []providers.Entry {
	marked := []providers.Entry{}
	for _, entry := range m.Entries {
		if entry.Marked {
			marked = append(marked, entry)
		}
	}
	return marked
}
