package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	leftPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		m.infoPanes[0].View(),
		m.panes[0].View(),
		m.panes[1].View(),
		m.infoPanes[1].View(),
	)

	rightPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		m.panes[2].View(),
		m.panes[3].View(),
	)

	topPanes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanes,
		rightPanes,
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		topPanes,
		m.statusBar.View(),
	)

	return view
}
