package simplelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/state"
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
		return m, refreshCmd(m.ParentIndex, m.Provider, m.state)

	case messages.RefreshResultMsg:
		m.Entries = providers.EntryNames(msg.Entries)
		m.updateTable()
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

	case messages.OpenMsg:
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) < 1 {
			return m, nil
		}

		selectedCol := selectedRow[0]
		if m.Active != nil && *m.Active == selectedCol {
			m.Active = nil
			selectedCol = ""

			if err := m.Provider.Deselect(); err != nil {
				return m, func() tea.Msg {
					return messages.ErrorMsg{Err: err}
				}
			}
		} else {
			m.Active = &selectedCol

			if err := m.Provider.Select(selectedCol); err != nil {
				return m, func() tea.Msg {
					return messages.ErrorMsg{Err: err}
				}
			}
		}

		var cmds []tea.Cmd

		for _, infoPaneIndex := range m.RefreshInfoPaneIndexes {
			cmd := func() tea.Msg {
				return messages.InfoPaneMsg{
					Index: infoPaneIndex,
					Msg:   messages.SetValueMsg{Value: selectedCol},
				}
			}
			cmds = append(cmds, cmd)
		}

		for _, paneIndex := range m.RefreshPaneIndexes {
			cmd := func() tea.Msg {
				switch paneIndex {

				case state.ProfilesPaneIndex:
					m.state.SetProfile(*m.Active)

				case state.BucketsPaneIndex:
					m.state.SetBucket(*m.Active)

				case state.RemotePaneIndex:
					// State is set on previous states

				default:
					// Unsupported
					return nil
				}

				return messages.PaneMsg{
					Index: paneIndex,
					Msg:   messages.RefreshMsg{},
				}
			}
			cmds = append(cmds, cmd)
		}

		return m, tea.Sequence(cmds...)
	}

	return m, nil
}

func refreshCmd(parentID int, provider providers.Provider, state *state.State) tea.Cmd {
	return func() tea.Msg {
		switch provider.(type) {
		case *providers.S3BucketProvider:
			if state.HasBucket() {
				args := []string{state.Bucket()}
				if state.HasRemotePath() {
					args = append(args, state.RemotePath())
				}

				if err := provider.Select(args...); err != nil {
					return messages.ErrorMsg{Err: err}
				}
			}

		case *providers.ProfilesProvider:
			if state.HasProfile() {
				if err := provider.Select(state.Profile()); err != nil {
					return messages.ErrorMsg{Err: err}
				}
			}

		default:
			// Unsupported
			return nil
		}

		entries, err := provider.Children()
		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		return messages.PaneMsg{
			Index: parentID,
			Msg:   messages.RefreshResultMsg{Entries: entries},
		}
	}
}

func (m *Model) updateTable() {
	rows := []table.Row{}
	for _, entry := range m.Entries {
		rows = append(rows, table.Row{entry})
	}
	m.table.SetRows(rows)
}
