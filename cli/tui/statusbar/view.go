package statusbar

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	h := help.New()
	h.Width = m.Width - 2
	h.ShortSeparator = " | "
	h.Styles.ShortDesc = h.Styles.ShortDesc.
		Foreground(render.ActivePalette.TextPrimary)

	statusBar := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Margin(0, 1).
		Render(h.ShortHelpView(m.KeyMap))

	return statusBar
}
