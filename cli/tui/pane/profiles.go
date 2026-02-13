package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	tea "github.com/charmbracelet/bubbletea"
)

type ProfilesPane struct {
	Model
}

func NewProfiles(id int, title string, opts ...PaneOption) ProfilesPane {
	return ProfilesPane{
		Model: New(id, title, opts...),
	}
}

func (m ProfilesPane) SetEntries(entries []providers.Entry) ProfilesPane {
	if list, ok := m.Child.(simplelist.Model); ok {
		m.Child = list.SetEntries(entries)
	}
	return m
}

func (m ProfilesPane) Clear() ProfilesPane {
	m.Model = m.Model.Clear()
	return m
}

func (m ProfilesPane) Update(msg tea.Msg) (ProfilesPane, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.OpenMsg:
		if list, ok := m.Child.(simplelist.Model); ok {
			if selection, ok := list.Selected(); ok {
				cmd = func() tea.Msg { return messages.ProfileSelectedMsg{Profile: selection} }
			}
		}

	default:
		m.Model, cmd = m.Model.Update(msg)
	}

	return m, cmd
}
