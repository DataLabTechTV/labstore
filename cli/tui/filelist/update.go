package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	UpdateTableMsg struct {
		Rows   []table.Row
		Cursor int
	}
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
		return m, refreshCmd(m.ParentIndex, m.Provider)

	case messages.RefreshResultMsg:
		m.Entries = msg.Entries
		m.updateTable(msg.Active)
		return m, nil

	case tea.FocusMsg:
		m.table.Focus()

	case tea.BlurMsg:
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

	case messages.LevelUpMsg:
		if err := m.Provider.Select(".."); err != nil {
			render.Error(err)
			return m, nil
		}

		cmd := func() tea.Msg {
			return messages.PaneMsg{
				Index: m.ParentIndex,
				Msg:   messages.RefreshMsg{},
			}
		}
		return m, cmd

	case messages.OpenMsg:
		path := m.table.SelectedRow()[2]

		if err := m.Provider.Select(path); err != nil {
			render.Error(err)
			return m, nil
		}

		cmd := func() tea.Msg {
			return messages.PaneMsg{
				Index: m.ParentIndex,
				Msg:   messages.RefreshMsg{},
			}
		}
		return m, cmd

	}

	return m, nil
}

func refreshCmd(parentIndex int, provider providers.Provider) tea.Cmd {
	return func() tea.Msg {
		entries, err := provider.Children()
		if err != nil {
			return messages.PaneMsg{
				Index: parentIndex,
				Msg:   messages.ErrorMsg{Err: err},
			}
		}

		var active *string
		if lastSelected, ok := provider.LastSelected(); ok && lastSelected != "" {
			active = &lastSelected
		}

		return messages.PaneMsg{
			Index: parentIndex,
			Msg:   messages.RefreshResultMsg{Entries: entries, Active: active},
		}
	}
}

func (m *Model) updateTable(active *string) {
	cursor := 0
	rows := []table.Row{}

	for i, entry := range m.Entries {
		var name, size string
		if entry.IsDir {
			name = entry.Name + "/"
			size = "-"
		} else {
			name = entry.Name
			size = render.NewSize(entry.Size).Format()
		}
		date := render.NewDate(entry.ModTime).Format()

		if active != nil && name == *active {
			cursor = i
		}

		rows = append(rows, table.Row{date, size, name})
	}

	m.table.SetRows(rows)
	m.table.SetCursor(cursor)
}
