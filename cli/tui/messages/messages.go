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

	FocusMsg struct{}
	BlurMsg  struct{}

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{ PaneIndex int }
	OpenMsg         struct{ PaneIndex int }

	SetValue struct{ Value string }
)
