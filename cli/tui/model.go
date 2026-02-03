package tui

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/infopane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/pane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/statusbar"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width       int
	height      int
	focusedPane int

	panes     []*pane.Model
	infoPanes []infopane.Model
	statusBar statusbar.Model
}

func New() Model {
	return Model{
		focusedPane: 1,
		panes: []*pane.Model{
			pane.New("[1] Buckets", pane.WithFocus()),
			pane.New("[2] Profiles"),
			pane.New("[3] Remote"),
			pane.New("[4] Local"),
		},
		infoPanes: []infopane.Model{
			infopane.New("Active Bucket", infopane.ValueNone),
			infopane.New("Active Profile", infopane.ValueNone),
		},
		statusBar: statusbar.New(DefaultKeyMap.All()),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
