package tui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/logger"
	"github.com/IllumiKnowLabs/labstore/client/types"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type progressMsg struct {
	current int
	total   int
}

type consoleMsg string

type ProgressBarModel struct {
	Bar progress.Model

	Console []string

	Progress chan types.Progress
	Message  chan string

	program *tea.Program
	width   int
	height  int

	Done chan struct{}
}

type ConsoleWriter struct {
	channel chan<- string
}

func NewProgressBarModel() (*ProgressBarModel, error) {
	m := &ProgressBarModel{
		Bar:     progress.New(progress.WithDefaultGradient()),
		Console: make([]string, 0),

		Progress: make(chan types.Progress, 10),
		Message:  make(chan string, 10),

		Done: make(chan struct{}),
	}

	m.program = tea.NewProgram(m)

	return m, nil
}

func (m *ProgressBarModel) Run() {
	output := ConsoleWriter{channel: m.Message}
	revert := logger.Temporary(output, logger.WithLevel(slog.LevelDebug))
	defer revert()
	defer close(m.Done)

	go func() {
		for msg := range m.Progress {
			m.program.Send(progressMsg{
				current: msg.Current,
				total:   msg.Total,
			})
		}
	}()

	go func() {
		for msg := range m.Message {
			m.program.Send(consoleMsg(msg))
		}
	}()

	if _, err := m.program.Run(); err != nil {
		slog.Error("progress bar", "err", err)
		return
	}

	slog.Debug("progress bar done")
}

func (m ProgressBarModel) Init() tea.Cmd {
	return m.Bar.Init()
}

func (m ProgressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			fmt.Println("SIGINT caught, upload canceled...")
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.Bar.Width = min(msg.Width, MaxWidth)
		return m, nil

	case progressMsg:
		pct := float64(msg.current) / float64(msg.total)
		cmd := m.Bar.SetPercent(pct)

		if m.Bar.Percent() >= 1.0 {
			return m, tea.Sequence(cmd, tea.Quit)
		}

		return m, cmd

	case consoleMsg:
		m.Console = append(m.Console, string(msg))
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.Bar.Update(msg)
		m.Bar = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m ProgressBarModel) View() string {
	wrappedLogs := lipgloss.NewStyle().
		Width(m.width).
		Render(strings.Join(m.Console, "\n"))

	barStyle := lipgloss.NewStyle().
		Margin(1, 0, 2, 0).
		Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wrappedLogs,
		barStyle(m.Bar.View()),
	)
}

func (w ConsoleWriter) Write(buf []byte) (n int, err error) {
	msg := strings.TrimRight(string(buf), "\n")
	w.channel <- msg
	return len(buf), nil
}
