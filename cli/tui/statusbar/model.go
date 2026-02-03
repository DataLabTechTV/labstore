package statusbar

import "github.com/charmbracelet/bubbles/key"

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
