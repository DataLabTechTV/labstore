package pane

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if m.Child != nil {
			m.Child, _ = m.Child.Update(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height})
		}
	}

	return m, nil
}
