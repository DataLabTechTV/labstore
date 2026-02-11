package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	tea "github.com/charmbracelet/bubbletea"
)

type BucketsPane struct {
	Model
}

func NewBuckets(id int, title string, opts ...PaneOption) BucketsPane {
	return BucketsPane{
		Model: New(id, title, opts...),
	}
}

func (m BucketsPane) SetEntries(entries []providers.Entry) BucketsPane {
	if list, ok := m.Child.(simplelist.Model); ok {
		m.Child = list.SetEntries(entries)
	}
	return m
}

func (m BucketsPane) Clear() BucketsPane {
	m.Model = m.Model.Clear()
	return m
}

func (m BucketsPane) Update(msg tea.Msg) (BucketsPane, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.OpenMsg:
		if list, ok := m.Child.(simplelist.Model); ok {
			if selection, ok := list.Selection(); ok {
				cmd = func() tea.Msg { return messages.BucketSelectedMsg{Bucket: selection} }
			}
		}

	default:
		m.Model, cmd = m.Model.Update(msg)
	}

	return m, cmd
}
