package multiprogress

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return m.Progress.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width

	case messages.UploadProgressMsg:
		m.Current[msg.FileIndex] = msg.Uploaded
		return m.updateProgress()

	case messages.DownloadProgressMsg:
		m.Current[msg.FileIndex] = msg.Downloaded
		return m.updateProgress()

	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateProgress() (Model, tea.Cmd) {
	var currentSum int64
	for _, current := range m.Current {
		currentSum += current
	}
	pct := float64(currentSum) / float64(m.Total)
	cmd := m.Progress.SetPercent(pct)
	return m, cmd
}
