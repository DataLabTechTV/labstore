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

	ErrorMsg struct{ Err error }

	RefreshMsg       struct{ Metadata map[string]string }
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
