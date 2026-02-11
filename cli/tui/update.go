package tui

import (
	"fmt"

	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/tui/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	statusBarHeight = 1
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// --- Left ---

		leftPaneWidth := int(1 / 3.0 * float64(m.width))

		m.infoPanes[0].Width = leftPaneWidth
		m.infoPanes[1].Width = leftPaneWidth
		remainingHeight := m.height - m.infoPanes[0].Height - m.infoPanes[1].Height - statusBarHeight

		cmds := make([]tea.Cmd, len(m.panes))

		m.panes[0], cmds[0] = m.panes[0].Update(tea.WindowSizeMsg{
			Width:  leftPaneWidth,
			Height: remainingHeight / 2,
		})
		remainingHeight -= m.panes[0].Height

		m.panes[1], cmds[1] = m.panes[1].Update(tea.WindowSizeMsg{
			Width:  leftPaneWidth,
			Height: remainingHeight,
		})

		// --- Right ---

		rightPaneWidth := m.width - leftPaneWidth

		m.panes[2], cmds[2] = m.panes[2].Update(tea.WindowSizeMsg{
			Width:  rightPaneWidth,
			Height: (m.height - statusBarHeight) / 2,
		})

		m.panes[3], cmds[3] = m.panes[3].Update(tea.WindowSizeMsg{
			Width:  rightPaneWidth,
			Height: m.height - m.panes[2].Height - statusBarHeight,
		})

		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch km := DefaultHomeKeyMap.(type) {
		case HomeKeyMap:
			switch {

			case key.Matches(msg, km.Refresh):
				var cmds []tea.Cmd

				m, cmd := m.sendToAll(messages.RefreshMsg{})
				cmds = append(cmds, cmd)

				alertCmd := func() tea.Msg {
					return messages.InfoMsg{
						Title:   "Refreshing",
						Content: "All panels are being refreshed",
					}
				}
				cmds = append(cmds, alertCmd)

				return m, tea.Batch(cmds...)

			case key.Matches(msg, km.Open):
				return m.sendToFocused(messages.OpenMsg{})

			case key.Matches(msg, km.NavUp):
				return m.sendToFocused(messages.LevelUpMsg{})

			case key.Matches(msg, km.Down):
				return m.sendToFocused(messages.MoveDownMsg{})

			case key.Matches(msg, km.Up):
				return m.sendToFocused(messages.MoveUpMsg{})

			case key.Matches(msg, km.End):
				return m.sendToFocused(messages.MoveToBottomMsg{})

			case key.Matches(msg, km.Home):
				return m.sendToFocused(messages.MoveToTopMsg{})

			case key.Matches(msg, km.PgDn):
				return m.sendToFocused(messages.PageDownMsg{})

			case key.Matches(msg, km.PgUp):
				return m.sendToFocused(messages.PageUpMsg{})

			case key.Matches(msg, km.Focus1):
				return m.paneFocus(0)

			case key.Matches(msg, km.Focus2):
				return m.paneFocus(1)

			case key.Matches(msg, km.Focus3):
				return m.paneFocus(2)

			case key.Matches(msg, km.Focus4):
				return m.paneFocus(3)

			case key.Matches(msg, km.Next):
				if m.focusedPaneIndex+1 >= len(m.panes) {
					return m.paneFocus(0)
				} else {
					return m.paneFocus(m.focusedPaneIndex + 1)
				}

			case key.Matches(msg, km.Previous):
				if m.focusedPaneIndex-1 < 0 {
					return m.paneFocus(len(m.panes) - 1)
				} else {
					return m.paneFocus(m.focusedPaneIndex - 1)
				}

			case key.Matches(msg, km.Quit):
				return m, tea.Quit
			}
		}

	case messages.InfoMsg:
		alert, cmd := alert.New(alert.AlertInfo, msg.Title, msg.Content)
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.WarnMsg:
		alert, cmd := alert.New(alert.AlertWarn, msg.Title, msg.Content)
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.ErrorMsg:
		alert, cmd := alert.New(alert.AlertError, fmt.Sprintf("%T", msg.Err), msg.Err.Error())
		m.alerts = append(m.alerts, alert)
		return m, cmd

	case messages.HideAlertMsg:
		alerts := []alert.Model{}
		for _, alert := range m.alerts {
			if alert.ID != msg.ID {
				alerts = append(alerts, alert)
			}
		}
		m.alerts = alerts
		return m, nil

	case messages.PaneMsg:
		var cmd tea.Cmd
		m.panes[msg.Index], cmd = m.panes[msg.Index].Update(msg.Msg)
		return m, cmd

	case messages.InfoPaneMsg:
		var cmd tea.Cmd
		m.infoPanes[msg.Index], cmd = m.infoPanes[msg.Index].Update(msg.Msg)

		switch msg.Index {

		case state.ProfileInfoIndex:
			m.State.SetProfile(msg.Msg.(messages.SetValueMsg).Value)

		case state.BucketInfoIndex:
			m.State.SetBucket(msg.Msg.(messages.SetValueMsg).Value)
		}

		return m, cmd
	}

	return m, nil
}

func (m Model) paneFocus(index int) (Model, tea.Cmd) {
	if index < 0 || index > len(m.panes) {
		return m, nil
	}

	var cmds []tea.Cmd
	m.focusedPaneIndex = index

	for i := range m.panes {
		var cmd tea.Cmd
		if i == index {
			m.panes[i], cmd = m.panes[i].Update(tea.FocusMsg{})
		} else {
			m.panes[i], cmd = m.panes[i].Update(tea.BlurMsg{})
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) sendToFocused(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.panes[m.focusedPaneIndex], cmd = m.panes[m.focusedPaneIndex].Update(msg)
	return m, cmd
}

func (m Model) sendToAll(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	for i := range m.panes {
		var cmd tea.Cmd
		m.panes[i], cmd = m.panes[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
