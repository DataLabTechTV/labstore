package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	tea "github.com/charmbracelet/bubbletea"
)

type LocalPane struct {
	Model
}

func NewLocal(id int, title string, opts ...PaneOption) LocalPane {
	return LocalPane{
		Model: New(id, title, opts...),
	}
}

func (m LocalPane) SetEntries(entries []providers.Entry, active *string) LocalPane {
	if list, ok := m.Child.(filelist.Model); ok {
		m.Child = list.SetEntries(entries, active)
	}
	return m
}

func (m LocalPane) Update(msg tea.Msg) (LocalPane, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}
