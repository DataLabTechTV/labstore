package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	markerChecked   = "☑"
	markerUnchecked = "☐"
)

type (
	UpdateTableMsg struct {
		Rows   []table.Row
		Cursor int
	}
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width < 1 || msg.Height < 1 {
			return m, nil
		}

		m.Width = msg.Width
		m.Height = msg.Height

		columns := make([]table.Column, 4)

		remainingWidth := m.Width

		selWidth := 2
		columns[0] = table.Column{Width: selWidth}
		remainingWidth -= selWidth

		modWidth := remainingWidth / 3
		remainingWidth -= modWidth
		columns[1] = table.Column{Title: "Modified", Width: modWidth}

		sizeWidth := remainingWidth / 4
		remainingWidth -= sizeWidth
		columns[2] = table.Column{Title: "Size", Width: sizeWidth}

		nameWidth := remainingWidth - m.hCellPad*len(columns)
		columns[3] = table.Column{Title: "Name", Width: nameWidth}

		m.table.SetWidth(m.Width)
		m.table.SetHeight(m.Height)
		m.table.SetColumns(columns)

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

	case messages.MarkMsg:
		m = m.Mark()
	}

	return m, nil
}

func (m *Model) updateTable(active *string) {
	cursor := 0
	rows := []table.Row{}

	for i, entry := range m.Entries {
		var name, size string
		if entry.IsDir {
			name = entry.Name
			size = "-"
		} else {
			name = entry.Name
			size = render.NewSize(entry.Size).Format()
		}
		date := render.NewDate(entry.ModTime).Format()

		if active != nil && entry.Path == *active {
			cursor = i
		}

		if entry.Marked {
			rows = append(rows, table.Row{markerChecked, date, size, name})
		} else {
			rows = append(rows, table.Row{markerUnchecked, date, size, name})
		}
	}

	if active == nil {
		cursor = m.table.Cursor()
	}
	m.table.SetRows(rows)
	m.table.SetCursor(cursor)
	m.table.MoveDown(0)
}
