package pane

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ID      int
	Title   string
	Focused bool
	Width   int
	Height  int
	Child   tea.Model
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

func WithSimpleList() PaneOption {
	return func(m *Model) {
		m.Child = simplelist.New()
	}
}

func WithFileList() PaneOption {
	return func(m *Model) {
		m.Child = filelist.New()
	}
}

func (m Model) SetFocused(focused bool) Model {
	m.Focused = focused
	return m
}

func (m Model) Clear() Model {
	switch child := m.Child.(type) {

	case simplelist.Model:
		m.Child = child.Clear()

	case filelist.Model:
		m.Child = child.Clear()
	}

	return m
}
