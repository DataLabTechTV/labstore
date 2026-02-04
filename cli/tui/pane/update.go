package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if m.Child != nil {
			// Not ideal to use magic numbers here, but it would be messy to access the boxStyle.
			m.Child, _ = m.Child.Update(tea.WindowSizeMsg{Width: msg.Width - 2, Height: msg.Height - 2})
		}

	case messages.FocusMsg:
		var cmd tea.Cmd
		m.Focused = true
		m.Child, cmd = m.Child.Update(messages.FocusMsg{})
		return m, cmd

	case messages.BlurMsg:
		var cmd tea.Cmd
		m.Focused = false
		m.Child, cmd = m.Child.Update(messages.BlurMsg{})
		return m, cmd

	case messages.MoveDownMsg, messages.MoveUpMsg:
		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg)
		return m, cmd

	case messages.FileListMsg:
		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg.Msg)
		return m, cmd

	case messages.SimpleListMsg:
		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg.Msg)
		return m, cmd
	}

	return m, nil
}
