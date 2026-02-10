package alert

import (
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var levelView string
	var borderColor lipgloss.Color
	switch m.Level {

	case AlertInfo:
		levelView = lipgloss.NewStyle().
			Foreground(render.ActivePalette.Success).
			Render("[INFO]")
		borderColor = render.ActivePalette.Success

	case AlertWarn:
		levelView = lipgloss.NewStyle().
			Foreground(render.ActivePalette.Warning).
			Render("[WARN]")
		borderColor = render.ActivePalette.Warning

	case AlertError:
		levelView = lipgloss.NewStyle().
			Foreground(render.ActivePalette.Error).
			Render("[ERROR]")
		borderColor = render.ActivePalette.Error

	default:
		levelView = lipgloss.NewStyle().
			Render("[Alert]")
		borderColor = render.ActivePalette.Border
	}

	border := lipgloss.ThickBorder()

	titleStyle := lipgloss.NewStyle().
		Border(border).
		BorderBottom(false).
		BorderForeground(borderColor).
		Width(m.Width).
		Render

	titleBorderBottom := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.MiddleLeft + strings.Repeat(border.Bottom, m.Width) + border.MiddleRight)

	titleView := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle(lipgloss.JoinHorizontal(lipgloss.Top, levelView, " ", m.Title)),
		titleBorderBottom,
	)

	boxStyle := lipgloss.NewStyle().
		Border(border, false, true, true, true).
		BorderForeground(borderColor).
		Width(m.Width).
		Height(1).
		Render

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		titleView,
		boxStyle(m.Message),
	)

	return view
}
