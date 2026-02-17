package pane

import (
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/types"
	"github.com/IllumiKnowLabs/labstore/server/helper"
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

func (m *LocalPane) SetEntries(entries []types.Entry, active *string) {
	if list, ok := m.Child.(filelist.Model); ok {
		list.SetEntries(entries, active)
		m.Child = list
	}
}

func (m *LocalPane) Clear() {
	m.Model.Clear()
}

func (m LocalPane) Update(msg tea.Msg) (LocalPane, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.OpenMsg:
		if list, ok := m.Child.(filelist.Model); ok {
			if dirname, ok := list.Selected(); ok {
				if dirname != ".." && !strings.HasSuffix(dirname, "/") {
					return m, nil
				}
				cmd = func() tea.Msg { return messages.LoadLocalMsg{Dirname: &dirname} }
			}
		}

	case messages.LevelUpMsg:
		cmd = func() tea.Msg { return messages.LoadLocalMsg{Dirname: helper.Ptr("..")} }

	default:
		m.Model, cmd = m.Model.Update(msg)
	}

	return m, cmd
}
