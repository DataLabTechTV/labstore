package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TUIModel struct {
	width  int
	height int
}

func NewTUIModel() TUIModel {
	return TUIModel{}
}

func (m TUIModel) Init() tea.Cmd {
	return nil
}
