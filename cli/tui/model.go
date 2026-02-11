package tui

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/infopane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/pane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/simplelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/state"
	"github.com/IllumiKnowLabs/labstore/cli/tui/statusbar"
	"github.com/IllumiKnowLabs/labstore/server/constants"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	State *state.State

	panes     []pane.Model
	infoPanes []infopane.Model
	statusBar statusbar.Model
	alerts    []alert.Model

	focusedPaneID int
	width         int
	height        int
}

func New() Model {
	var globalState state.State

	s3BucketProvider := providers.NewS3BucketProvider()
	profilesProvider := providers.NewProfilesProvider()
	s3FSProvider := providers.NewS3FSProvider()
	fsProvider := providers.NewFSProvider()

	bucketList := simplelist.New(
		&globalState,
		state.BucketsPaneIndex,
		s3BucketProvider,
		simplelist.WithRefreshInfoPaneIndexes([]int{state.BucketInfoIndex}),
		simplelist.WithRefreshPaneIndexes([]int{state.RemotePaneIndex}),
	)
	profileList := simplelist.New(
		&globalState,
		state.ProfilesPaneIndex,
		profilesProvider,
		simplelist.WithRefreshInfoPaneIndexes([]int{state.ProfileInfoIndex}),
		simplelist.WithRefreshPaneIndexes([]int{state.BucketsPaneIndex}),
	)
	s3FileList := filelist.New(&globalState, state.RemotePaneIndex, s3FSProvider)
	fsFileList := filelist.New(&globalState, state.LocalPaneIndex, fsProvider)

	bucketPane := pane.New(
		1, "Buckets",
		pane.WithFocus(),
		pane.WithChild(bucketList),
	)

	profilesPane := pane.New(
		2, "Profiles",
		pane.WithChild(profileList),
	)

	remotePane := pane.New(
		3, "Remote",
		pane.WithChild(s3FileList),
		pane.WithProvider(s3FSProvider),
	)

	localPane := pane.New(
		4, "Local",
		pane.WithChild(fsFileList),
		pane.WithProvider(fsProvider),
	)

	bucketInfo := infopane.New("Active Bucket", infopane.ValueNone)
	profileInfo := infopane.New("Active Profile", infopane.ValueNone)

	return Model{
		State:         &globalState,
		panes:         []pane.Model{bucketPane, profilesPane, remotePane, localPane},
		infoPanes:     []infopane.Model{bucketInfo, profileInfo},
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
		if paneCmd := pane.Init(); paneCmd != nil {
			cmd := func() tea.Msg {
				msg := paneCmd()
				if msg == nil {
					return nil
				}
				return messages.PaneMsg{Index: i, Msg: msg}
			}
			cmds = append(cmds, cmd)
		}
	}

	cmds = append(cmds, m.statusBar.Init())
	cmds = append(cmds, tea.SetWindowTitle(fmt.Sprintf("%s TUI", constants.Name)))

	return tea.Batch(cmds...)
}
