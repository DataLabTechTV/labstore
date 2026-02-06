package messages

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	PaneMsg struct {
		ID  int
		Msg tea.Msg
	}

	InfoPaneMsg struct {
		ID  int
		Msg tea.Msg
	}

	RefreshMsg       struct{}
	RefreshResultMsg struct {
		Entries []providers.Entry
		Active  *string
		Err     error
	}

	SetValueMsg struct{ Value string }

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{}
	OpenMsg         struct{}
)
