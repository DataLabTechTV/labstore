package filelist

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Provider providers.Provider
	Entries  []providers.Entry
	Width    int
	Height   int
}

func New(provider providers.Provider) Model {
	return Model{
		Provider: provider,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return messages.FileListMsg{Msg: messages.RefreshMsg{}}
	}
}
