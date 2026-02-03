package filelist

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	RefreshMsg  struct{}
	MoveUpMsg   struct{}
	MoveDownMsg struct{}
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case RefreshMsg:
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		entries, err := m.Provider.List(ctx)
		if err != nil {
			render.Error(err)
		}
		m.Entries = entries

	case MoveUpMsg:
	case MoveDownMsg:
	}

	return m, nil
}
