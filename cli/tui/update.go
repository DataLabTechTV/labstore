package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	statusBarHeight = 1
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// --- Left ---

		leftPaneWidth := int(1 / 3.0 * float64(m.width))

		m.infoPanes[0].Width = leftPaneWidth
		m.infoPanes[1].Width = leftPaneWidth
		remainingHeight := m.height - m.infoPanes[0].Height - m.infoPanes[1].Height - statusBarHeight

		cmds := make([]tea.Cmd, len(m.panes))

		m.panes[0], cmds[0] = m.panes[0].Update(tea.WindowSizeMsg{
			Width:  leftPaneWidth,
			Height: remainingHeight / 2,
		})
		remainingHeight -= m.panes[0].Height

		m.panes[1], cmds[1] = m.panes[1].Update(tea.WindowSizeMsg{
			Width:  leftPaneWidth,
			Height: remainingHeight,
		})

		// --- Right ---

		rightPaneWidth := m.width - leftPaneWidth

		m.panes[2], cmds[2] = m.panes[2].Update(tea.WindowSizeMsg{
			Width:  rightPaneWidth,
			Height: (m.height - statusBarHeight) / 2,
		})

		m.panes[3], cmds[3] = m.panes[3].Update(tea.WindowSizeMsg{
			Width:  rightPaneWidth,
			Height: m.height - m.panes[2].Height - statusBarHeight,
		})

		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch km := DefaultHomeKeyMap.(type) {
		case HomeKeyMap:
			switch {
			// case key.Matches(msg, km.Down):
			// 	activePane := m.panes[m.focusedPane]
			case key.Matches(msg, km.Focus1):
				m.paneFocus(1)
				return m, nil

			case key.Matches(msg, km.Focus2):
				m.paneFocus(2)
				return m, nil

			case key.Matches(msg, km.Focus3):
				m.paneFocus(3)
				return m, nil

			case key.Matches(msg, km.Focus4):
				m.paneFocus(4)
				return m, nil

			case key.Matches(msg, km.Next):
				if m.focusedPane >= len(m.panes) {
					m.paneFocus(1)
				} else {
					m.paneFocus(m.focusedPane + 1)
				}

			case key.Matches(msg, km.Previous):
				if m.focusedPane <= 1 {
					m.paneFocus(4)
				} else {
					m.paneFocus(m.focusedPane - 1)
				}

			case key.Matches(msg, km.Quit):
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m *Model) paneFocus(n int) {
	for i, pane := range m.panes {
		if i == n-1 {
			pane.Focused = true
		} else {
			pane.Focused = false
		}
	}

	m.focusedPane = n
}
