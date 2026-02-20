package confirm

import tea "github.com/charmbracelet/bubbletea"

func (m Model[T, U]) Init() tea.Cmd {
	return nil
}

func (m Model[T, U]) Update(msg tea.Msg) (Model[T, U], tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}
