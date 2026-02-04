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

	focusedPaneID int
	width         int
	height        int
}

func New() Model {
	bucketProvider := simplelist.New(1, providers.NewS3BucketProvider())
	profilesProvider := simplelist.New(2, providers.NewProfilesProvider())
	s3FSProvider := filelist.New(3, providers.NewS3FSProvider())
	fsProvider := filelist.New(4, providers.NewFSProvider("."))

	bucketPane := pane.New(1, "Buckets", pane.WithFocus(), pane.WithChild(bucketProvider))
	profilesPane := pane.New(2, "Profiles", pane.WithChild(profilesProvider))
	remotePane := pane.New(3, "Remote", pane.WithChild(s3FSProvider))
	localPane := pane.New(4, "Local", pane.WithChild(fsProvider))

	activeBucketsInfo := infopane.New("Active Bucket", infopane.ValueNone)
	activeProfileInfo := infopane.New("Active Profile", infopane.ValueNone)

	return Model{
		panes:         []pane.Model{bucketPane, profilesPane, remotePane, localPane},
		infoPanes:     []infopane.Model{activeBucketsInfo, activeProfileInfo},
		statusBar:     statusbar.New(DefaultHomeKeyMap.HelpKeys()),
		focusedPaneID: 1,
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
