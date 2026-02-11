package pane

import tea "github.com/charmbracelet/bubbletea"

type LocalPane struct {
	Model
}

func NewLocal(id int, title string, opts ...PaneOption) LocalPane {
	return LocalPane{
		Model: New(id, title, opts...),
	}
}

func (m LocalPane) Update(msg tea.Msg) (LocalPane, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}
