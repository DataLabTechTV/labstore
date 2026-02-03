package filelist

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/bubbles/table"
)

func (m Model) View() string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entries, err := m.Provider.List(ctx, ".")
	if err != nil {
		render.Error(err)
	}

	const totalPadding = 3 * 2 // 3 columns, cell padding of 1+1 (left+right)
	modifiedWidth := m.Width / 3
	sizeWidth := (m.Width - modifiedWidth) / 4
	nameWidth := m.Width - modifiedWidth - sizeWidth - totalPadding

	columns := []table.Column{
		{Title: "Modified", Width: modifiedWidth},
		{Title: "Size", Width: sizeWidth},
		{Title: "Name", Width: nameWidth},
	}

	rows := []table.Row{}
	for _, entry := range entries {
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

	tableStyle := table.DefaultStyles()

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
