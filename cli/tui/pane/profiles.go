package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	"github.com/IllumiKnowLabs/labstore/cli/types"
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

func (m *ProfilesPane) SetEntries(entries []types.Entry) {
	if list, ok := m.Child.(simplelist.Model); ok {
		list.SetEntries(entries)
		m.Child = list
	}
}

func (m *ProfilesPane) EntryNames() []string {
	if list, ok := m.Child.(simplelist.Model); ok {
		return list.GetEntryNames()
	}
	return nil
}

func (m *ProfilesPane) Clear() {
	m.Model.Clear()
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
