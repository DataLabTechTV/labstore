package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ID       int
	Title    string
	Focused  bool
	Width    int
	Height   int
	Child    tea.Model
	Provider providers.Provider
}

type PaneOption func(m *Model)

func New(id int, title string, opts ...PaneOption) Model {
	m := Model{
		ID:      id,
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

func WithProvider(provider providers.Provider) PaneOption {
	return func(m *Model) {
		m.Provider = provider
	}
}

func (m Model) Init() tea.Cmd {
	if childCmd := m.Child.Init(); childCmd != nil {
		return childCmd
	}
	return nil
}
