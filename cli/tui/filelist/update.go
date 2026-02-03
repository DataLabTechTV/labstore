package filelist

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	MoveUpMsg   struct{}
	MoveDownMsg struct{}
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case MoveUpMsg:
	case MoveDownMsg:
	}

	return m, nil
}
