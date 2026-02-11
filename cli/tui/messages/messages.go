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

	AlertInfoMsg  struct{ Title, Content string }
	AlertWarnMsg  struct{ Title, Content string }
	AlertErrorMsg struct{ Err error }
	AlertHideMsg  struct{ ID string }

	RefreshMsg       struct{}
	RefreshResultMsg struct {
		Entries []providers.Entry
		Active  *string
	}

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
