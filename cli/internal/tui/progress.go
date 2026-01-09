package tui

import (
	"log/slog"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/logger"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ConsoleLines = 5
)

type ProgressMsg struct {
	current, total int
}

type ConsoleMsg string

type ProgressUI struct {
	console viewport.Model
	bar     progress.Model
	logs    []string
}

type ConsoleWriter struct {
	program *tea.Program
}

func NewProgressBar() (s3.ProgressCallback, error) {
	program := tea.NewProgram(
		&ProgressUI{
			console: viewport.New(0, ConsoleLines),
			bar:     progress.New(progress.WithDefaultGradient()),
			logs:    make([]string, 0),
		},
	)

	progressCh := make(chan ProgressMsg, 10)

	callback := func(current, total int) {
		progressCh <- ProgressMsg{current: current, total: total}
	}

	output := ConsoleWriter{program: program}
	revert := logger.Temporary(output)

	go func() {
		defer revert()

		if _, err := program.Run(); err != nil {
			slog.Error("progress bar", "err", err)
			return
		}

		slog.Debug("end progress bar")
	}()

	go func() {
		for msg := range progressCh {
			program.Send(msg)
		}
	}()

	return callback, nil
}

func (ui ProgressUI) Init() tea.Cmd {
	return nil
}

func (ui ProgressUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ui.console.Width = min(msg.Width, MaxWidth)
		ui.bar.Width = min(msg.Width, MaxWidth)
		return ui, nil

	case ProgressMsg:
		if ui.bar.Percent() >= 1.0 {
			return ui, tea.Quit
		}

		pct := float64(msg.current) / float64(msg.total)
		cmd := ui.bar.SetPercent(pct)
		return ui, cmd

	case ConsoleMsg:
		ui.logs = append(ui.logs, string(msg))
		ui.console.SetContent(strings.Join(ui.logs, "\n"))
		ui.console.GotoBottom()
		return ui, nil

	case progress.FrameMsg:
		progressModel, cmd := ui.bar.Update(msg)
		ui.bar = progressModel.(progress.Model)
		return ui, cmd

	default:
		return ui, nil
	}
}

func (ui ProgressUI) View() string {
	consoleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Render

	barStyle := lipgloss.NewStyle().
		Margin(1, 0).
		Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		consoleStyle(ui.console.View()),
		barStyle(ui.bar.View()),
		// HelpStyle("Press any key to quit"),
	)
}

func (w ConsoleWriter) Write(buf []byte) (n int, err error) {
	msg := strings.TrimRight(string(buf), "\n")
	w.program.Send(ConsoleMsg(msg))
	return len(buf), nil
}
