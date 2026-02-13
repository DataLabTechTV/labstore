package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		columns := []table.Column{{Width: m.Width - m.hCellPad}}
		m.table.SetColumns(columns)
		m.table.SetWidth(m.Width)
		m.table.SetHeight(m.Height)

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

func (m *Model) updateTable() {
	rows := []table.Row{}
	for _, entry := range m.Entries {
		rows = append(rows, table.Row{entry})
	}
	m.table.SetRows(rows)
	m.table.GotoTop()
}
