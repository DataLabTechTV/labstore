package pane

import (
	"fmt"
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Width < 2 || m.Height < 2 {
		return ""
	}

	border := lipgloss.RoundedBorder()

	boxStyle := lipgloss.NewStyle().
		Width(m.Width-border.GetLeftSize()-border.GetRightSize()).
		Height(m.Height-border.GetTopSize()-border.GetBottomSize()).
		Border(border, false, true, true)

	titleStyle := lipgloss.NewStyle()
	borderStyle := lipgloss.NewStyle()

	if m.Focused {
		titleStyle = titleStyle.Foreground(render.ActivePalette.Accent)
		borderStyle = borderStyle.Foreground(render.ActivePalette.SurfaceHover)
		boxStyle = boxStyle.BorderForeground(render.ActivePalette.SurfaceHover)
	} else {
		titleStyle = titleStyle.Foreground(render.ActivePalette.AccentMuted)
		borderStyle = borderStyle.Foreground(render.ActivePalette.Surface)
		boxStyle = boxStyle.BorderForeground(render.ActivePalette.Surface)
	}

	var topBorder string
	if m.Width >= len(m.Title)+4 {
		topLeftBorder := borderStyle.Render(border.TopLeft + border.Top)
		topRightBorder := borderStyle.Render(strings.Repeat(border.Top, m.Width-len(m.Title)-3) + border.TopRight)
		topBorder = topLeftBorder + titleStyle.Render(fmt.Sprintf("[%d] %s", m.ID, m.Title)) + topRightBorder
	} else {
		boxStyle = boxStyle.BorderTop(true)
	}

	var childView string
	if m.Child != nil {
		childView = m.Child.View()
	} else {
		childView = ""
	}

	pane := lipgloss.JoinVertical(
		lipgloss.Top,
		topBorder,
		boxStyle.Render(childView),
	)

	return pane
}
