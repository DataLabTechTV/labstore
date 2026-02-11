package tui

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/charmbracelet/bubbles/key"
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

		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch km := DefaultHomeKeyMap.(type) {
		case HomeKeyMap:
			switch {

			case key.Matches(msg, km.Refresh):
				var cmds []tea.Cmd

				m, cmd := m.SendToAllPanes(messages.RefreshMsg{})
				cmds = append(cmds, cmd)

				alertCmd := func() tea.Msg {
					return messages.AlertInfoMsg{
						Title:   "Refreshing",
						Content: "All panels are being refreshed",
					}
				}
				cmds = append(cmds, alertCmd)

				return m, tea.Batch(cmds...)

			case key.Matches(msg, km.Open):
				return m.SendToFocusedPane(messages.OpenMsg{})

			case key.Matches(msg, km.NavUp):
				return m.SendToFocusedPane(messages.LevelUpMsg{})

			case key.Matches(msg, km.Down):
				return m.SendToFocusedPane(messages.MoveDownMsg{})

			case key.Matches(msg, km.Up):
				return m.SendToFocusedPane(messages.MoveUpMsg{})

			case key.Matches(msg, km.End):
				return m.SendToFocusedPane(messages.MoveToBottomMsg{})

			case key.Matches(msg, km.Home):
				return m.SendToFocusedPane(messages.MoveToTopMsg{})

			case key.Matches(msg, km.PgDn):
				return m.SendToFocusedPane(messages.PageDownMsg{})

			case key.Matches(msg, km.PgUp):
				return m.SendToFocusedPane(messages.PageUpMsg{})

			case key.Matches(msg, km.Focus1):
				var cmd tea.Cmd
				m = m.SetFocusedPane(FocusBuckets)
				m.bucketsPane, cmd = m.bucketsPane.Update(messages.RefreshMsg{})
				return m, cmd

			case key.Matches(msg, km.Focus2):
				m = m.SetFocusedPane(FocusProfiles)

			case key.Matches(msg, km.Focus3):
				m = m.SetFocusedPane(FocusRemote)

			case key.Matches(msg, km.Focus4):
				m = m.SetFocusedPane(FocusLocal)

			case key.Matches(msg, km.Next):
				m = m.SetFocusedPane((m.focusedPane + 1) % 4)

			case key.Matches(msg, km.Previous):
				m = m.SetFocusedPane((m.focusedPane + 4 - 1) % 4)

			case key.Matches(msg, km.Quit):
				return m, tea.Quit
			}
		}

	case messages.LoadProfilesMsg:
		return m, func() tea.Msg {
			entries, err := m.profilesProvider.Children()
			if err != nil {
				return messages.AlertErrorMsg{Err: err}
			}
			return messages.ProfilesLoadedMsg{Entries: entries}
		}

	case messages.ProfilesLoadedMsg:
		var cmd tea.Cmd
		m.profilesPane = m.profilesPane.SetEntries(msg.Entries)
		m.profilesPane, cmd = m.profilesPane.Update(msg)
		return m, cmd

	case messages.ProfileSelectedMsg:
		if m.profileInfoPane.Value == msg.Profile {
			m.profileInfoPane = m.profileInfoPane.Clear()

			m.bucketsPane = m.bucketsPane.Clear()
			m.bucketInfoPane = m.bucketInfoPane.Clear()
			m.remotePane = m.remotePane.Clear()

			m.AppState.UnsetProfile()
			m.AppState.UnsetBucket()
			m.AppState.UnsetRemotePath()
		} else {
			m.profileInfoPane.Value = msg.Profile
			m.AppState.SetProfile(msg.Profile)
			return m, func() tea.Msg { return messages.LoadBucketsMsg(msg) }
		}

	case messages.LoadBucketsMsg:
		return m, func() tea.Msg {
			if err := m.s3BucketProvider.Select(msg.Profile); err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			entries, err := m.s3BucketProvider.Children()
			if err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			return messages.BucketsLoadedMsg{Entries: entries}
		}

	case messages.BucketsLoadedMsg:
		var cmd tea.Cmd
		m.bucketsPane = m.bucketsPane.SetEntries(msg.Entries)
		m.bucketsPane, cmd = m.bucketsPane.Update(msg)
		return m, cmd

	case messages.BucketSelectedMsg:
		if m.bucketInfoPane.Value == msg.Bucket {
			m.bucketInfoPane = m.bucketInfoPane.Clear()

			m.remotePane = m.remotePane.Clear()

			m.AppState.UnsetBucket()
			m.AppState.UnsetRemotePath()
		} else {
			m.bucketInfoPane.Value = msg.Bucket
			m.AppState.SetBucket(msg.Bucket)
			return m, func() tea.Msg { return messages.LoadRemoteMsg{} }
		}

	case messages.LoadRemoteMsg:
		if msg.Dirname != nil {
			if *msg.Dirname == ".." && m.AppState.IsRemotePathRoot() {
				return m, nil
			}
			m.AppState.CDRemotePath(*msg.Dirname)
		}

		return m, func() tea.Msg {
			var args []string

			if !m.AppState.HasProfile() {
				return messages.AlertErrorMsg{Err: &errs.ErrProfileNotSelected{}}
			}
			args = append(args, m.AppState.Profile())

			if !m.AppState.HasBucket() {
				return messages.AlertErrorMsg{Err: &errs.ErrBucketNotSelected{}}
			}
			args = append(args, m.AppState.Bucket())

			if m.AppState.HasRemotePath() {
				args = append(args, m.AppState.RemotePath())
			}

			if err := m.s3FSProvider.Select(args...); err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			entries, err := m.s3FSProvider.Children()
			if err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			var active *string
			if lastSel, ok := m.s3BucketProvider.LastSelected(); ok {
				active = &lastSel
			}

			return messages.RemoteLoadedMsg{Entries: entries, Active: active}
		}

	case messages.RemoteLoadedMsg:
		var cmd tea.Cmd
		m.remotePane = m.remotePane.SetEntries(msg.Entries, msg.Active)
		m.remotePane, cmd = m.remotePane.Update(msg)
		return m, cmd

	case messages.LoadLocalMsg:
		if msg.Dirname != nil {
			if *msg.Dirname == ".." && m.AppState.IsLocalPathRoot() {
				return m, nil
			}
			m.AppState.CDLocalPath(*msg.Dirname)
		}

		return m, func() tea.Msg {
			if !m.AppState.HasLocalPath() {
				return func() tea.Msg { return messages.AlertErrorMsg{Err: &errs.ErrLocalPathNotSet{}} }
			}

			if err := m.fsProvider.Select(m.AppState.LocalPath()); err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			entries, err := m.fsProvider.Children()
			if err != nil {
				return messages.AlertErrorMsg{Err: err}
			}

			var active *string
			if lastSel, ok := m.fsProvider.LastSelected(); ok {
				active = &lastSel
			}

			return messages.LocalLoadedMsg{Entries: entries, Active: active}
		}

	case messages.LocalLoadedMsg:
		var cmd tea.Cmd
		m.localPane = m.localPane.SetEntries(msg.Entries, msg.Active)
		m.localPane, cmd = m.localPane.Update(msg)
		return m, cmd

	case messages.AlertInfoMsg:
		alert, cmd := alert.New(alert.AlertInfo, msg.Title, msg.Content)
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.AlertWarnMsg:
		alert, cmd := alert.New(alert.AlertWarn, msg.Title, msg.Content)
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.AlertErrorMsg:
		alert, cmd := alert.New(alert.AlertError, fmt.Sprintf("%T", msg.Err), msg.Err.Error())
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.AlertHideMsg:
		alerts := []alert.Model{}
		for _, alert := range m.alerts {
			if alert.ID != msg.ID {
				alerts = append(alerts, alert)
			}
		}
		m.alerts = alerts
		return m, nil
	}

	return m, nil
}

func (m Model) SetFocusedPane(focusedPane FocusedPane) Model {
	m.focusedPane = focusedPane
	m.bucketsPane.Model = m.bucketsPane.SetFocused(focusedPane == FocusBuckets)
	m.profilesPane.Model = m.profilesPane.SetFocused(focusedPane == FocusProfiles)
	m.remotePane.Model = m.remotePane.SetFocused(focusedPane == FocusRemote)
	m.localPane.Model = m.localPane.SetFocused(focusedPane == FocusLocal)
	return m
}

func (m Model) SendToFocusedPane(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) SendToAllPanes(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	m.bucketsPane, cmd = m.bucketsPane.Update(msg)
	cmds = append(cmds, cmd)

	m.profilesPane, cmd = m.profilesPane.Update(msg)
	cmds = append(cmds, cmd)

	m.remotePane, cmd = m.remotePane.Update(msg)
	cmds = append(cmds, cmd)

	m.localPane, cmd = m.localPane.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
