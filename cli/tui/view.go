package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
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

	var alerts []string
	var totalHeight int
	for _, alert := range m.alerts {
		alertView := alert.View()
		totalHeight += strings.Count(alertView, "\n")
		alerts = append(alerts, alertView)
	}
	alertsView := lipgloss.JoinVertical(
		lipgloss.Top,
		alerts...,
	)

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		topPanes,
		m.statusBar.View(),
	)

	view := overlay.Composite(alertsView, mainView, overlay.Right, overlay.Top, -2, 1)

	return view
}
