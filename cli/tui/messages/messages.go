package messages

import tea "github.com/charmbracelet/bubbletea"

type (
	PaneMsg struct {
		Index int
		Msg   tea.Msg
	}

	FileListMsg struct {
		Msg tea.Msg
	}

	SimpleListMsg struct {
		Msg tea.Msg
	}
)
