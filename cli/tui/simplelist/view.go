package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/bubbles/table"
)

func (m Model) View() string {
	columns := make([]table.Column, 3)

	tableStyle := table.DefaultStyles()
	totalPadding := tableStyle.Cell.GetHorizontalPadding() * len(columns)

	modifiedWidth := m.Width / 3
	columns[0] = table.Column{Title: "Modified", Width: modifiedWidth}

	sizeWidth := (m.Width - modifiedWidth) / 4
	columns[1] = table.Column{Title: "Size", Width: sizeWidth}

	nameWidth := m.Width - modifiedWidth - sizeWidth - totalPadding
	columns[2] = table.Column{Title: "Name", Width: nameWidth}

	rows := []table.Row{}
	for _, entry := range m.Entries {
		row := table.Row{
			render.NewDate(entry.ModTime).Format(),
			render.NewSize(entry.Size).Format(),
			entry.Name,
		}
		rows = append(rows, row)
	}

	fileTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithWidth(m.Width),
		table.WithHeight(m.Height),
	)

	tableStyle.Header = tableStyle.Header.
		Foreground(render.ActivePalette.SurfaceAlt).
		Bold(true)

	tableStyle.Selected = tableStyle.Selected.
		Foreground(render.ActivePalette.Accent).
		Background(render.ActivePalette.SurfaceAlt).
		Bold(false)

	fileTable.SetStyles(tableStyle)

	return fileTable.View()
}
