package pane

import (
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/server/helper"
	tea "github.com/charmbracelet/bubbletea"
)

type RemotePane struct {
	Model
}

func NewRemote(id int, title string, opts ...PaneOption) RemotePane {
	return RemotePane{
		Model: New(id, title, opts...),
	}
}

func (m *RemotePane) SetEntries(entries []providers.Entry, active *string) {
	if list, ok := m.Child.(filelist.Model); ok {
		list.SetEntries(entries, active)
		m.Child = list
	}
}

func (m *RemotePane) Clear() {
	m.Model.Clear()
}

func (m RemotePane) Update(msg tea.Msg) (RemotePane, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.OpenMsg:
		if list, ok := m.Child.(filelist.Model); ok {
			if dirname, ok := list.Selected(); ok {
				if dirname != ".." && !strings.HasSuffix(dirname, "/") {
					return m, nil
				}
				cmd = func() tea.Msg { return messages.LoadRemoteMsg{Dirname: &dirname} }
			}
		}

	case messages.LevelUpMsg:
		cmd = func() tea.Msg { return messages.LoadRemoteMsg{Dirname: helper.Ptr("..")} }

	default:
		m.Model, cmd = m.Model.Update(msg)
	}

	return m, cmd
}
