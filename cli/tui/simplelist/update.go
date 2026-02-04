package simplelist

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		columns := []table.Column{{Width: m.Width - m.hCellPad}}
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
		m.Entries = providers.EntryNames(entries)

		rows := []table.Row{}
		for _, entry := range m.Entries {
			rows = append(rows, table.Row{entry})
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

	case messages.OpenMsg:
		selectedProfile := m.table.SelectedRow()[0]

		cmd := func() tea.Msg {
			return messages.InfoPaneMsg{
				Index: msg.PaneIndex - 1,
				Msg:   messages.SetValue{Value: selectedProfile},
			}
		}

		return m, cmd
	}

	return m, nil
}
