package infopane

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Width < 2 || m.Height < 2 {
		return ""
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(render.ActivePalette.AccentMuted).
		Render

	var valueColor lipgloss.Color
	if m.Value == ValueNone {
		valueColor = render.ActivePalette.TextMuted
	} else {
		valueColor = render.ActivePalette.Accent
	}

	valueStyle := lipgloss.NewStyle().
		Foreground(valueColor).
		Render

	var displayValue string
	if m.Value == "" {
		displayValue = "<none>"
	} else {
		displayValue = m.Value
	}

	infoPane := lipgloss.NewStyle().
		Width(m.Width - 2).
		Height(m.Height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(render.ActivePalette.Surface).
		Render(labelStyle(m.Label) + ": " + valueStyle(displayValue))

	return infoPane
}
