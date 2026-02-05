package messages

import tea "github.com/charmbracelet/bubbletea"

type (
	PaneMsg struct {
		Index int
		Msg   tea.Msg
	}

	InfoPaneMsg struct {
		Index int
		Msg   tea.Msg
	}

	FileListMsg struct {
		Msg tea.Msg
	}

	SimpleListMsg struct {
		Msg tea.Msg
	}

	RefreshMsg struct{}

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{}
	OpenMsg         struct{}

	SetTitle struct{ Title string }
	SetValue struct{ Value string }
)
