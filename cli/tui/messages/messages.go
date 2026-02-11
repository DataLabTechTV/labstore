package messages

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	PaneMsg struct {
		Index int
		Msg   tea.Msg
	}

	InfoPaneMsg struct {
		Index int
		Msg   tea.Msg
	}

	InfoMsg  struct{ Title, Content string }
	WarnMsg  struct{ Title, Content string }
	ErrorMsg struct{ Err error }

	RefreshMsg       struct{}
	RefreshResultMsg struct {
		Entries []providers.Entry
		Active  *string
	}

	HideAlertMsg struct{ ID string }

	SelectItemMsg struct{ Item string }
	SetValueMsg   struct{ Value string }

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{}
	OpenMsg         struct{}
)
