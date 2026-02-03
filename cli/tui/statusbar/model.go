package statusbar

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	KeyMap []key.Binding
	Width  int
	Height int
}

func New(keyMap []key.Binding) Model {
	return Model{
		KeyMap: keyMap,
		Height: 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
