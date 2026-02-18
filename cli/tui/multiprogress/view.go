package multiprogress

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		Padding(2, 4).
		Width(m.Width).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	m.Progress.Width = m.Width - boxStyle.GetHorizontalFrameSize()

	return boxStyle.Render(m.Progress.View())
}
