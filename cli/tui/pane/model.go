package pane

import tea "github.com/charmbracelet/bubbletea"

type PaneOption func(m *Model)

type Model struct {
	Title   string
	Width   int
	Height  int
	Focused bool
	ViewFn  func(width, height int) string
}

func New(title string, opts ...PaneOption) *Model {
	m := &Model{
		Title:   title,
		Focused: false,
		ViewFn: func(width, height int) string {
			return ""
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithFocus() PaneOption {
	return func(m *Model) {
		m.Focused = true
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
