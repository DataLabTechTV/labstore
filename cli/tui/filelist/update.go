package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/state"
	"github.com/IllumiKnowLabs/labstore/server/helper"
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
		return m, refreshCmd(m.ParentIndex, m.Provider, m.state)

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
		if !m.state.HasLocalPath() {
			return m, nil
		}

		cmd := func() tea.Msg {
			m.state.CDLocalPath("..")
			return messages.PaneMsg{
				Index: m.ParentIndex,
				Msg:   messages.RefreshMsg{},
			}
		}
		return m, cmd

	case messages.OpenMsg:
		dirname := m.table.SelectedRow()[2]

		if dirname == ".." {
			return m, nil
		}

		if !helper.IsDir(dirname) {
			return m, nil
		}

		switch m.ParentIndex {
		case state.LocalPaneIndex:
			m.state.CDLocalPath(dirname)

		case state.RemotePaneIndex:
			m.state.CDRemotePath(dirname)

		default:
			// Unsupported
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

func refreshCmd(parentIndex int, provider providers.Provider, state *state.State) tea.Cmd {
	return func() tea.Msg {
		switch provider.(type) {
		case *providers.FSProvider:
			if state.HasLocalPath() {
				if err := provider.Select(state.LocalPath()); err != nil {
					return messages.ErrorMsg{Err: err}
				}
			}

		case *providers.S3FSProvider:
			if state.HasProfile() && state.HasBucket() {
				args := []string{state.Profile(), state.Bucket()}
				if state.HasRemotePath() {
					args = append(args, state.RemotePath())
				}

				if err := provider.Select(args...); err != nil {
					return messages.ErrorMsg{Err: err}
				}
			}

		default:
			// Unsupported
			return nil
		}

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

		rows = append(rows, table.Row{date, size, name})
	}

	m.table.SetRows(rows)
	m.table.SetCursor(cursor)
}
