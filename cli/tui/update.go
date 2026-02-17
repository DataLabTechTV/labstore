package tui

import (
	"fmt"
	"os"

	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/cli/tui/filelist"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/multiprogress"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	statusBarHeight = 1
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return messages.LoadProfilesMsg{} },
		func() tea.Msg { return messages.LoadLocalMsg{} },
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.HandleWindowSize(msg)

	case tea.KeyMsg:
		return m.HandleKey(msg)

	case progress.FrameMsg:
		return m.HandleProgressFrameMsg(msg)

	case messages.LoadProfilesMsg:
		return m.HandleLoadProfiles(msg)

	case messages.ProfilesLoadedMsg:
		return m.HandleProfilesLoaded(msg)

	case messages.ProfilesFailedMsg:
		return m.HandleProfilesFailed(msg)

	case messages.ProfileSelectedMsg:
		return m.HandleProfileSelected(msg)

	case messages.LoadBucketsMsg:
		return m.HandleLoadBuckets(msg)

	case messages.BucketsLoadedMsg:
		return m.HandleBucketsLoaded(msg)

	case messages.BucketsFailedMsg:
		return m.HandleBucketsFailed(msg)

	case messages.BucketSelectedMsg:
		return m.HandleBucketSelected(msg)

	case messages.LoadRemoteMsg:
		return m.HandleLoadRemote(msg)

	case messages.RemoteLoadedMsg:
		return m.HandleRemoteLoaded(msg)

	case messages.RemoteFailedMsg:
		return m.HandleRemoteFailed(msg)

	case messages.LoadLocalMsg:
		return m.HandleLoadLocal(msg)

	case messages.LocalLoadedMsg:
		return m.HandleLocalLoaded(msg)

	case messages.LocalFailedMsg:
		return m.HandleLocalFailed(msg)

	case messages.RefreshAllMsg:
		return m.HandleRefreshAll(msg)

	case messages.StartUploadMsg:
		return m.HandleStartUpload(msg)

	case messages.UploadProgressMsg:
		return m.HandleUploadProgress(msg)

	case messages.UploadDoneMsg:
		return m.HandleUploadDone(msg)

	case messages.AlertInfoMsg:
		return m.HandleAlertInfo(msg)

	case messages.AlertWarnMsg:
		return m.HandleAlertWarn(msg)

	case messages.AlertErrorMsg:
		return m.HandleAlertError(msg)

	case messages.AlertHideMsg:
		return m.HandleAlertHide(msg)
	}

	return m, nil
}

func (m Model) HandleWindowSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	// --- Left ---

	leftPaneWidth := int(1 / 3.0 * float64(m.width))

	m.bucketInfoPane.Width = leftPaneWidth
	m.profileInfoPane.Width = leftPaneWidth
	remainingHeight := m.height - m.bucketInfoPane.Height - m.profileInfoPane.Height - statusBarHeight

	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.bucketsPane, cmd = m.bucketsPane.Update(tea.WindowSizeMsg{
		Width:  leftPaneWidth,
		Height: remainingHeight / 2,
	})
	cmds = append(cmds, cmd)
	remainingHeight -= m.bucketsPane.Height

	m.profilesPane, cmd = m.profilesPane.Update(tea.WindowSizeMsg{
		Width:  leftPaneWidth,
		Height: remainingHeight,
	})
	cmds = append(cmds, cmd)

	// --- Right ---

	rightPaneWidth := m.width - leftPaneWidth

	m.remotePane, cmd = m.remotePane.Update(tea.WindowSizeMsg{
		Width:  rightPaneWidth,
		Height: (m.height - statusBarHeight) / 2,
	})
	cmds = append(cmds, cmd)

	m.localPane, cmd = m.localPane.Update(tea.WindowSizeMsg{
		Width:  rightPaneWidth,
		Height: m.height - m.remotePane.Height - statusBarHeight,
	})
	cmds = append(cmds, cmd)

	if m.multiProgress != nil {
		m.multiProgress.Update(tea.WindowSizeMsg{
			Width:  m.width - m.width/4,
			Height: m.height - m.height/4,
		})
	}

	return m, tea.Batch(cmds...)
}

func (m Model) HandleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch km := DefaultHomeKeyMap.(type) {
	case HomeKeyMap:
		switch {

		case key.Matches(msg, km.Put):
			return m, func() tea.Msg { return messages.StartUploadMsg{} }

		case key.Matches(msg, km.Refresh):
			var cmds []tea.Cmd

			refreshCmd := func() tea.Msg { return messages.RefreshAllMsg{} }
			cmds = append(cmds, refreshCmd)

			alertCmd := func() tea.Msg {
				return messages.AlertInfoMsg{
					Title:   "Refreshing",
					Content: "All panels are being refreshed",
				}
			}
			cmds = append(cmds, alertCmd)

			return m, tea.Batch(cmds...)

		case key.Matches(msg, km.Open):
			return m.sendToFocusedPane(messages.OpenMsg{})

		case key.Matches(msg, km.NavUp):
			return m.sendToFocusedPane(messages.LevelUpMsg{})

		case key.Matches(msg, km.Select):
			return m.sendToFocusedPane(messages.MarkMsg{})

		case key.Matches(msg, km.Down):
			return m.sendToFocusedPane(messages.MoveDownMsg{})

		case key.Matches(msg, km.Up):
			return m.sendToFocusedPane(messages.MoveUpMsg{})

		case key.Matches(msg, km.End):
			return m.sendToFocusedPane(messages.MoveToBottomMsg{})

		case key.Matches(msg, km.Home):
			return m.sendToFocusedPane(messages.MoveToTopMsg{})

		case key.Matches(msg, km.PgDn):
			return m.sendToFocusedPane(messages.PageDownMsg{})

		case key.Matches(msg, km.PgUp):
			return m.sendToFocusedPane(messages.PageUpMsg{})

		case key.Matches(msg, km.Focus1):
			m.setFocusedPane(FocusBuckets)

		case key.Matches(msg, km.Focus2):
			m.setFocusedPane(FocusProfiles)

		case key.Matches(msg, km.Focus3):
			m.setFocusedPane(FocusRemote)

		case key.Matches(msg, km.Focus4):
			m.setFocusedPane(FocusLocal)

		case key.Matches(msg, km.Next):
			m.setFocusedPane((m.focusedPane + 1) % 4)

		case key.Matches(msg, km.Previous):
			m.setFocusedPane((m.focusedPane + 4 - 1) % 4)

		case key.Matches(msg, km.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) HandleProgressFrameMsg(msg progress.FrameMsg) (Model, tea.Cmd) {
	if m.multiProgress != nil {
		multiProgress, cmd := m.multiProgress.Update(msg)
		m.multiProgress = &multiProgress
		return m, cmd
	}
	return m, nil
}

func (m Model) HandleLoadProfiles(msg messages.LoadProfilesMsg) (Model, tea.Cmd) {
	return m, m.loadProfilesCmd()
}

func (m Model) HandleProfilesLoaded(msg messages.ProfilesLoadedMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.profilesPane.SetEntries(msg.Entries)
	m.profilesPane, cmd = m.profilesPane.Update(msg)
	return m, cmd
}

func (m Model) HandleProfilesFailed(msg messages.ProfilesFailedMsg) (Model, tea.Cmd) {
	return m, func() tea.Msg { return messages.AlertErrorMsg{Err: msg.Err} }
}

func (m Model) HandleProfileSelected(msg messages.ProfileSelectedMsg) (Model, tea.Cmd) {
	if m.profileInfoPane.Value == msg.Profile {
		m.deselectProfile()
		return m, nil
	}

	m.profileInfoPane.Value = msg.Profile
	m.AppState.SetProfile(msg.Profile)
	return m, func() tea.Msg { return messages.LoadBucketsMsg{} }
}

func (m Model) HandleLoadBuckets(msg messages.LoadBucketsMsg) (Model, tea.Cmd) {
	if !m.AppState.HasProfile() {
		m.deselectProfile()
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: &errs.ErrNoProfileSelected{}} }
	}

	return m, m.loadBucketsCmd()
}

func (m Model) HandleBucketsLoaded(msg messages.BucketsLoadedMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.bucketsPane.SetEntries(msg.Entries)
	m.bucketsPane, cmd = m.bucketsPane.Update(msg)
	return m, cmd
}

func (m Model) HandleBucketsFailed(msg messages.BucketsFailedMsg) (Model, tea.Cmd) {
	m.resetBuckets()
	return m, func() tea.Msg { return messages.AlertErrorMsg{Err: msg.Err} }
}

func (m Model) HandleBucketSelected(msg messages.BucketSelectedMsg) (Model, tea.Cmd) {
	m.resetRemote()

	if m.bucketInfoPane.Value == msg.Bucket {
		m.deselectBucket()
		m.s3BucketProvider.Deselect()
		m.s3FSProvider.Deselect()
	} else {
		m.bucketInfoPane.Value = msg.Bucket
		m.AppState.SetBucket(msg.Bucket)
		return m, func() tea.Msg { return messages.LoadRemoteMsg{} }
	}

	return m, nil
}

func (m Model) HandleLoadRemote(msg messages.LoadRemoteMsg) (Model, tea.Cmd) {
	if msg.Dirname != nil {
		if *msg.Dirname == ".." && m.AppState.IsRemotePathRoot() {
			return m, nil
		}
		m.AppState.CDRemotePath(*msg.Dirname)
	}

	var args []string

	if !m.AppState.HasProfile() {
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: &errs.ErrNoProfileSelected{}} }
	}
	args = append(args, m.AppState.Profile())

	if !m.AppState.HasBucket() {
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: &errs.ErrNoBucketSelected{}} }
	}
	args = append(args, m.AppState.Bucket())

	if m.AppState.HasRemotePath() {
		args = append(args, m.AppState.RemotePath())
	}

	if err := m.s3FSProvider.Select(args...); err != nil {
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: err} }
	}

	return m, m.loadRemoteCmd()
}

func (m Model) HandleRemoteLoaded(msg messages.RemoteLoadedMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.remotePane.SetEntries(msg.Entries, msg.Active)
	m.remotePane, cmd = m.remotePane.Update(msg)
	return m, cmd
}

func (m Model) HandleRemoteFailed(msg messages.RemoteFailedMsg) (Model, tea.Cmd) {
	return m, func() tea.Msg { return messages.AlertErrorMsg{Err: msg.Err} }
}

func (m Model) HandleLoadLocal(msg messages.LoadLocalMsg) (Model, tea.Cmd) {
	if msg.Dirname == nil {
		path, err := os.Getwd()
		if err != nil {
			return m, func() tea.Msg { return messages.AlertErrorMsg{Err: err} }
		}
		m.AppState.SetLocalPath(path)
	} else {
		if *msg.Dirname == ".." && m.AppState.IsLocalPathRoot() {
			return m, nil
		}
		m.AppState.CDLocalPath(*msg.Dirname)
	}

	if !m.AppState.HasLocalPath() {
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: &errs.ErrLocalPathNotSet{}} }
	}

	if err := m.fsProvider.Select(m.AppState.LocalPath()); err != nil {
		return m, func() tea.Msg { return messages.AlertErrorMsg{Err: err} }
	}

	return m, m.loadLocalCmd()
}

func (m Model) HandleLocalLoaded(msg messages.LocalLoadedMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.localPane.SetEntries(msg.Entries, msg.Active)
	m.localPane, cmd = m.localPane.Update(msg)
	return m, cmd
}

func (m Model) HandleLocalFailed(msg messages.LocalFailedMsg) (Model, tea.Cmd) {
	return m, func() tea.Msg { return messages.AlertErrorMsg{Err: msg.Err} }
}

func (m Model) HandleRefreshAll(msg messages.RefreshAllMsg) (Model, tea.Cmd) {
	m.bucketInfoPane.Clear()
	m.profileInfoPane.Clear()
	m.bucketsPane.Clear()
	m.profilesPane.Clear()
	m.remotePane.Clear()
	m.localPane.Clear()

	return m, tea.Batch(
		m.loadProfilesCmd(),
		m.loadBucketsCmd(),
		m.loadRemoteCmd(),
		m.loadLocalCmd(),
	)
}

func (m Model) HandleStartUpload(msg messages.StartUploadMsg) (Model, tea.Cmd) {
	fileList, ok := m.localPane.Child.(filelist.Model)
	if !ok {
		return m, func() tea.Msg {
			return messages.AlertErrorMsg{
				Title:   "Upload Failed",
				Content: "Local pane does not contain a file list",
			}
		}
	}

	srcs := fileList.Marked()
	if len(srcs) < 1 {
		return m, func() tea.Msg {
			return messages.AlertErrorMsg{
				Title:   "Upload Failed",
				Content: "No selected source files",
			}
		}
	}

	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return messages.AlertInfoMsg{
			Title:   "Upload Started",
			Content: fmt.Sprintf("Uploading %d files to %s/%s", len(srcs), m.s3FSProvider.Bucket, m.s3FSProvider.Key),
		}
	}
	cmds = append(cmds, cmd)

	m.multiProgress = multiprogress.New(len(srcs))
	cmds = append(cmds, m.multiProgress.Init())

	for _, src := range srcs {
		stat, err := os.Stat(src.Path)
		if err != nil {
			return m, func() tea.Msg { return messages.AlertErrorMsg{Err: err} }
		}
		m.multiProgress.Total += stat.Size()
	}

	m.multiProgress.UploadProgressCh = m.s3FSProvider.Upload(m.fsProvider.Path, srcs...)
	cmds = append(cmds, m.waitForUploadProgressCmd())

	return m, tea.Sequence(cmds...)
}

func (m Model) waitForUploadProgressCmd() tea.Cmd {
	return func() tea.Msg {
		progressMsg, ok := <-m.multiProgress.UploadProgressCh
		if !ok {
			return messages.UploadDoneMsg{}
		}

		if progressMsg.Err != nil {
			return messages.AlertErrorMsg{Err: progressMsg.Err}
		}

		return progressMsg
	}
}

func (m Model) HandleUploadProgress(msg messages.UploadProgressMsg) (Model, tea.Cmd) {
	multiProgress, cmd := m.multiProgress.Update(msg)
	m.multiProgress = &multiProgress
	return m, tea.Sequence(cmd, m.waitForUploadProgressCmd())
}

func (m Model) HandleUploadDone(msg messages.UploadDoneMsg) (Model, tea.Cmd) {
	return m, func() tea.Msg {
		return messages.AlertInfoMsg{
			Title:   "Upload Success",
			Content: fmt.Sprintf("Successfully uploaded %d files", msg.FileCount),
		}
	}
}

func (m Model) HandleAlertInfo(msg messages.AlertInfoMsg) (Model, tea.Cmd) {
	alert, cmd := alert.New(alert.AlertInfo, msg.Title, msg.Content)
	m.alerts = append(m.alerts, alert)
	return m, cmd
}

func (m Model) HandleAlertWarn(msg messages.AlertWarnMsg) (Model, tea.Cmd) {
	alert, cmd := alert.New(alert.AlertWarn, msg.Title, msg.Content)
	m.alerts = append(m.alerts, alert)
	return m, cmd
}

func (m Model) HandleAlertError(msg messages.AlertErrorMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if msg.Err != nil {
		alert, cmd := alert.New(alert.AlertError, fmt.Sprintf("%T", msg.Err), msg.Err.Error())
		cmds = append(cmds, cmd)
		m.alerts = append(m.alerts, alert)
	}

	if msg.Title != "" && msg.Content != "" {
		alert, cmd := alert.New(alert.AlertError, msg.Title, msg.Content)
		cmds = append(cmds, cmd)
		m.alerts = append(m.alerts, alert)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) HandleAlertHide(msg messages.AlertHideMsg) (Model, tea.Cmd) {
	alerts := []alert.Model{}
	for _, alert := range m.alerts {
		if alert.ID != msg.ID {
			alerts = append(alerts, alert)
		}
	}
	m.alerts = alerts
	return m, nil
}

func (m *Model) setFocusedPane(focusedPane FocusedPane) {
	m.focusedPane = focusedPane
	m.bucketsPane.Model = m.bucketsPane.SetFocused(focusedPane == FocusBuckets)
	m.profilesPane.Model = m.profilesPane.SetFocused(focusedPane == FocusProfiles)
	m.remotePane.Model = m.remotePane.SetFocused(focusedPane == FocusRemote)
	m.localPane.Model = m.localPane.SetFocused(focusedPane == FocusLocal)
}

func (m Model) sendToFocusedPane(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.focusedPane {

	case FocusBuckets:
		m.bucketsPane, cmd = m.bucketsPane.Update(msg)

	case FocusProfiles:
		m.profilesPane, cmd = m.profilesPane.Update(msg)

	case FocusRemote:
		m.remotePane, cmd = m.remotePane.Update(msg)

	case FocusLocal:
		m.localPane, cmd = m.localPane.Update(msg)
	}

	return m, cmd
}

func (m Model) loadProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		entries, err := m.profilesProvider.Children()
		if err != nil {
			return messages.ProfilesFailedMsg{Err: err}
		}
		return messages.ProfilesLoadedMsg{Entries: entries}
	}
}

func (m Model) loadBucketsCmd() tea.Cmd {
	return func() tea.Msg {
		if err := m.s3BucketProvider.Select(m.AppState.Profile()); err != nil {
			return messages.BucketsFailedMsg{Err: err}
		}

		entries, err := m.s3BucketProvider.Children()
		if err != nil {
			return messages.BucketsFailedMsg{Err: err}
		}

		return messages.BucketsLoadedMsg{Entries: entries}
	}
}

func (m Model) loadRemoteCmd() tea.Cmd {
	return func() tea.Msg {
		entries, err := m.s3FSProvider.Children()
		if err != nil {
			return messages.RemoteFailedMsg{Err: err}
		}

		var active *string
		if lastSel, ok := m.s3FSProvider.LastSelected(); ok {
			active = &lastSel
		}

		return messages.RemoteLoadedMsg{Entries: entries, Active: active}
	}
}

func (m Model) loadLocalCmd() tea.Cmd {
	return func() tea.Msg {
		entries, err := m.fsProvider.Children()
		if err != nil {
			return messages.LocalFailedMsg{Err: err}
		}

		var active *string
		if lastSel, ok := m.fsProvider.LastSelected(); ok {
			active = &lastSel
		}

		return messages.LocalLoadedMsg{Entries: entries, Active: active}
	}
}

func (m *Model) deselectProfile() {
	m.AppState.UnsetProfile()
	m.profilesProvider.Deselect()
	m.profileInfoPane.Clear()
	m.resetBuckets()
}

func (m *Model) deselectBucket() {
	m.AppState.UnsetBucket()
	m.s3BucketProvider.Deselect()
	m.bucketInfoPane.Clear()
	m.resetRemote()
}

func (m *Model) resetBuckets() {
	m.deselectBucket()
	m.bucketsPane.Clear()
}

func (m *Model) resetRemote() {
	m.AppState.UnsetRemotePath()
	m.s3FSProvider.Deselect()
	m.remotePane.Clear()
}
