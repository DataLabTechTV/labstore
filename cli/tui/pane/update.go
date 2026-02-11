package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	if childCmd := m.Child.Init(); childCmd != nil {
		return childCmd
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if m.Child != nil {
			// Not ideal to use magic numbers here, but it would be messy to access the boxStyle.
			m.Child, _ = m.Child.Update(tea.WindowSizeMsg{Width: msg.Width - 2, Height: msg.Height - 2})
		}

	case tea.FocusMsg:
		var cmd tea.Cmd
		m.Focused = true
		m.Child, cmd = m.Child.Update(tea.FocusMsg{})
		return m, cmd

	case tea.BlurMsg:
		var cmd tea.Cmd
		m.Focused = false
		m.Child, cmd = m.Child.Update(tea.BlurMsg{})
		return m, cmd

	case messages.MoveDownMsg, messages.MoveUpMsg,
		messages.MoveToBottomMsg, messages.MoveToTopMsg,
		messages.PageDownMsg, messages.PageUpMsg,
		messages.LevelUpMsg:

		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg)
		return m, cmd

	case messages.OpenMsg:
		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg)
		return m, cmd

	case messages.ProfilesLoadedMsg:
		var cmd tea.Cmd
		if m.Child != nil {
			m.Child, cmd = m.Child.Update(msg)
		}
		return m, cmd

	case messages.RefreshMsg:
		var cmd tea.Cmd
		if m.Child != nil {
			m.Child, cmd = m.Child.Update(msg)
		}
		return m, cmd

		// case messages.RefreshResultMsg:
		// 	if m.Provider != nil {
		// 		if selected := m.Provider.Selected(); selected != "" {
		// 			m.Title = selected
		// 		}
		// 	}

		// 	child, cmd := m.Child.Update(msg)
		// 	m.Child = child
		// 	return m, cmd
	}

	return m, nil
}
