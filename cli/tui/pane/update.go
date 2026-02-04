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

	case messages.FileListMsg:
		var cmd tea.Cmd
		m.Child, cmd = m.Child.Update(msg.Msg)
		return m, cmd
	}

	return m, nil
}
