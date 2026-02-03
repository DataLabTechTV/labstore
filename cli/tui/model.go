package tui

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/infopane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/pane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/statusbar"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	panes     []pane.Model
	infoPanes []infopane.Model
	statusBar statusbar.Model

	focusedPane int
	width       int
	height      int
}

func New() Model {
	localFileList := filelist.New(providers.NewFSProvider())

	return Model{
		panes: []pane.Model{
			pane.New("[1] Buckets", pane.WithFocus()),
			pane.New("[2] Profiles"),
			pane.New("[3] Remote"),
			pane.New("[4] Local", pane.WithChild(localFileList)),
		},
		infoPanes: []infopane.Model{
			infopane.New("Active Bucket", infopane.ValueNone),
			infopane.New("Active Profile", infopane.ValueNone),
		},
		statusBar: statusbar.New(DefaultHomeKeyMap.All()),

		focusedPane: 1,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, infoPane := range m.infoPanes {
		cmds = append(cmds, infoPane.Init())
	}

	for _, pane := range m.panes {
		cmds = append(cmds, pane.Init())
	}

	cmds = append(cmds, m.statusBar.Init())

	return tea.Batch(cmds...)
}
