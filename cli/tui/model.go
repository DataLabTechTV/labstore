package tui

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/infopane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/pane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/statusbar"
	"github.com/IllumiKnowLabs/labstore/server/constants"
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
	bucketProvider := simplelist.New(providers.NewS3BucketProvider())
	profilesProvider := simplelist.New(providers.NewProfilesProvider())
	s3FSProvider := simplelist.New(providers.NewS3FSProvider())
	fsProvider := filelist.New(providers.NewFSProvider("."))

	return Model{
		panes: []pane.Model{
			pane.New("[1] Buckets", pane.WithFocus(), pane.WithChild(bucketProvider)),
			pane.New("[2] Profiles", pane.WithChild(profilesProvider)),
			pane.New("[3] Remote", pane.WithChild(s3FSProvider)),
			pane.New("[4] Local", pane.WithChild(fsProvider)),
		},
		infoPanes: []infopane.Model{
			infopane.New("Active Bucket", infopane.ValueNone),
			infopane.New("Active Profile", infopane.ValueNone),
		},
		statusBar: statusbar.New(DefaultHomeKeyMap.All()),
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, infoPane := range m.infoPanes {
		cmds = append(cmds, infoPane.Init())
	}

	for i, pane := range m.panes {
		paneCmd := pane.Init()

		if paneCmd == nil {
			continue
		}

		cmd := func() tea.Msg {
			msg := paneCmd()
			if msg == nil {
				return nil
			}
			return messages.PaneMsg{Index: i, Msg: msg}
		}
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, m.statusBar.Init())
	cmds = append(cmds, tea.SetWindowTitle(fmt.Sprintf("%s TUI", constants.Name)))

	return tea.Batch(cmds...)
}
