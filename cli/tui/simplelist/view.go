package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/bubbles/table"
)

func (m Model) View() string {
	if len(m.Entries) < 1 {
		return ""
	}

	tableStyle := table.DefaultStyles()
	columns := []table.Column{{Width: m.Width - tableStyle.Cell.GetHorizontalPadding()}}

	rows := []table.Row{}
	for _, entry := range m.Entries {
		rows = append(rows, table.Row{entry})
	}

	simpleTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithWidth(m.Width),
		table.WithHeight(m.Height),
	)

	tableStyle.Selected = tableStyle.Selected.
		Foreground(render.ActivePalette.Accent).
		Background(render.ActivePalette.SurfaceAlt).
		Bold(false)

	simpleTable.SetStyles(tableStyle)

	return simpleTable.View()
}
