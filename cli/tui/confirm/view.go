package confirm

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/lipgloss"
)

const (
	cancelBtnText  = "Cancel"
	confirmBtnText = "Confirm"
)

func (m Model[T, U]) View() string {
	boxStyle := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Padding(1).
		Border(lipgloss.ThickBorder()).
		Render

	btnStyle := lipgloss.NewStyle().
		Padding(1).
		Margin(0, 1).
		Background(render.ActivePalette.Surface).
		Foreground(render.ActivePalette.TextPrimary).
		Render

	btnActiveStyle := lipgloss.NewStyle().
		Padding(1).
		Margin(0, 1).
		Background(render.ActivePalette.SurfaceHover).
		Foreground(render.ActivePalette.TextInverted).
		Render

	var cancelBtnView string
	if !m.HoverConfirm {
		cancelBtnView = btnActiveStyle(cancelBtnText)
	} else {
		cancelBtnView = btnStyle(cancelBtnText)
	}

	var confirmBtnView string
	if m.HoverConfirm {
		confirmBtnView = btnActiveStyle(confirmBtnText)
	} else {
		confirmBtnView = btnStyle(confirmBtnText)
	}

	btnView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cancelBtnView,
		confirmBtnView,
	)

	contentView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.Prompt,
		btnView,
	)

	return boxStyle(contentView)
}
