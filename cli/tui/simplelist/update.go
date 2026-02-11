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
		value := selectedRow[0]

		if m.Active != nil && *m.Active == value {
			m.Active = nil
			value = ""
		} else {
			m.Active = &value
		}

		switch m.ParentIndex {

		case state.ProfilesPaneIndex:
			if m.Active == nil {
				m.state.UnsetProfile()
			} else {
				m.state.SetProfile(*m.Active)
			}

		case state.BucketsPaneIndex:
			if m.Active == nil {
				m.state.UnsetBucket()
			} else {
				m.state.SetBucket(*m.Active)
			}

		case state.RemotePaneIndex:
			// State is set on previous states
		}

		var cmds []tea.Cmd

		for _, infoPaneIndex := range m.RefreshInfoPaneIndexes {
			cmd := func() tea.Msg {
				return messages.InfoPaneMsg{
					Index: infoPaneIndex,
					Msg:   messages.SetValueMsg{Value: value},
				}
			}
			cmds = append(cmds, cmd)
		}

		for _, paneIndex := range m.RefreshPaneIndexes {
			cmd := func() tea.Msg {
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

func refreshCmd(parentIndex int, provider providers.Provider, globalState *state.State) tea.Cmd {
	return func() tea.Msg {
		switch parentIndex {

		case state.ProfilesPaneIndex:
			if globalState.HasProfile() {
				if err := provider.Select(globalState.Profile()); err != nil {
					return messages.ErrorMsg{Err: err}
				}
			}

		case state.BucketsPaneIndex:
			if globalState.HasProfile() {
				if err := provider.Select(globalState.Profile()); err != nil {
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
			Index: parentIndex,
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
