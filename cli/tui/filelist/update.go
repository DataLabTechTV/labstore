package filelist

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case messages.RefreshMsg:
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		entries, err := m.Provider.List(ctx)
		if err != nil {
			render.Error(err)
			return m, nil
		}
		m.Entries = entries

	case messages.MoveDownMsg:
	case messages.MoveUpMsg:
	}

	return m, nil
}
