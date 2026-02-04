package filelist

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		columns := make([]table.Column, 3)

		modifiedWidth := m.Width / 3
		columns[0] = table.Column{Title: "Modified", Width: modifiedWidth}

		sizeWidth := (m.Width - modifiedWidth) / 4
		columns[1] = table.Column{Title: "Size", Width: sizeWidth}

		nameWidth := m.Width - modifiedWidth - sizeWidth - m.hCellPad*len(columns)
		columns[2] = table.Column{Title: "Name", Width: nameWidth}

		m.table.SetColumns(columns)
		m.table.SetWidth(m.Width)
		m.table.SetHeight(m.Height)

	case messages.RefreshMsg:
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		entries, err := m.Provider.List(ctx)
		if err != nil {
			render.Error(err)
			return m, nil
		}
		m.Entries = entries

		rows := []table.Row{}
		for _, entry := range m.Entries {
			var name, size string
			if entry.IsDir {
				name = entry.Name + "/"
				size = "-"
			} else {
				name = entry.Name
				size = render.NewSize(entry.Size).Format()
			}
			date := render.NewDate(entry.ModTime).Format()

			rows = append(rows, table.Row{date, size, name})
		}
		m.table.SetRows(rows)

	case messages.FocusMsg:
		m.table.Focus()

	case messages.BlurMsg:
		m.table.Blur()

	case messages.MoveDownMsg:
		if last := len(m.table.Rows()) - 1; m.table.Cursor() == last {
			m.table.GotoTop()
		} else {
			m.table.MoveDown(1)
		}

	case messages.MoveUpMsg:
		if m.table.Cursor() == 0 {
			m.table.GotoBottom()
		} else {
			m.table.MoveUp(1)
		}

	case messages.MoveToBottomMsg:
		m.table.GotoBottom()

	case messages.MoveToTopMsg:
		m.table.GotoTop()

	case messages.PageDownMsg:
		m.table.MoveDown(10)

	case messages.PageUpMsg:
		m.table.MoveUp(10)
	}

	return m, nil
}
