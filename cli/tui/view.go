package tui

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/server/constants"
	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

func (m Model) View() string {
	leftPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		m.bucketInfoPane.View(),
		m.bucketsPane.View(),
		m.profilesPane.View(),
		m.profileInfoPane.View(),
	)

	rightPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		m.remotePane.View(),
		m.localPane.View(),
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

	var version string
	if constants.GitTag == constants.Unknown {
		version = fmt.Sprintf("v%s", constants.Version)
	} else {
		version = constants.GitTag
	}

	titleBarView := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(m.width).
		Foreground(render.ActivePalette.TextMuted).
		Bold(true).
		PaddingRight(1).
		Render(fmt.Sprintf("%s %s, by %s, %s", constants.Name, version, constants.Author, constants.GitRepo))

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		titleBarView,
		topPanes,
		m.statusBar.View(),
	)

	view := overlay.Composite(alertsView, mainView, overlay.Left, overlay.Top, m.width-alert.DefaultWidth-2, 1)

	if m.multiProgress != nil {
		progressBarView := m.multiProgress.View()
		view = overlay.Composite(progressBarView, view, overlay.Center, overlay.Center, 0, 0)
	}

	return view
}
