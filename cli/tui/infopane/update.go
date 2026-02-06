package infopane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case messages.SetValueMsg:
		if msg.Value == "" {
			m.Value = ValueNone
		} else {
			m.Value = msg.Value
		}
	}

	return m, nil
}
