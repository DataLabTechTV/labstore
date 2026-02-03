package statusbar

import (
	"github.com/IllumiKnowLabs/labstore/cli/render"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	helpMsg := help.New()
	helpMsg.ShortSeparator = " | "
	helpMsg.Width = m.Width - 2
	helpMsg.Styles.ShortDesc = helpMsg.Styles.ShortDesc.
		Foreground(render.ActivePalette.TextPrimary)

	statusBar := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Margin(0, 1).
		Render(helpMsg.ShortHelpView(m.KeyMap))

	return statusBar
}
