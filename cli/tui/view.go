package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	statusBarHeight = 1
)

func (m TUIModel) View() string {
	leftPaneWidth := int(1 / 3.0 * float64(m.width))

	bottomPaneHeight := int((float64(m.height-statusBarHeight-2) / 2.0))
	topPaneHeight := m.height - bottomPaneHeight - statusBarHeight - 2

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())

	infoPane := boxStyle.
		Width(leftPaneWidth).
		Height(topPaneHeight).
		Render()

	credentialsPane := boxStyle.
		Width(leftPaneWidth).
		Height(bottomPaneHeight).
		Render()

	localPane := boxStyle.
		Width(m.width - leftPaneWidth - 2).
		Height(topPaneHeight).
		Render()

	remotePane := boxStyle.
		Width(m.width - leftPaneWidth - 2).
		Height(bottomPaneHeight).
		Render()

	statusBar := lipgloss.NewStyle().
		Width(m.width).
		Height(statusBarHeight).
		Margin(0, 1).
		Foreground(lipgloss.Color("003")).
		Render("Profile: P | PUT: p | GET: g | HEAD: h")

	topPanes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		infoPane,
		localPane,
	)

	bottomPanes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		credentialsPane,
		remotePane,
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		topPanes,
		bottomPanes,
		statusBar,
	)

	return view
}
