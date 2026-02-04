package pane

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Title   string
	Focused bool
	Width   int
	Height  int
	Child   tea.Model
}

type PaneOption func(m *Model)

func New(title string, opts ...PaneOption) Model {
	m := Model{
		Title:   title,
		Focused: false,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithFocus() PaneOption {
	return func(m *Model) {
		m.Focused = true
	}
}

func WithChild(child tea.Model) PaneOption {
	return func(m *Model) {
		m.Child = child
	}
}

func (m Model) Init() tea.Cmd {
	if m.Child != nil {
		childCmd := m.Child.Init()
		return func() tea.Msg {
			msg := childCmd()
			if msg == nil {
				return nil
			}
			return msg
		}
	}

	return nil
}
