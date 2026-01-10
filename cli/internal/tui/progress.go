package tui

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/logger"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressMsg struct {
	current, total int
}

type ConsoleMsg string

type ProgressUI struct {
	Width   int
	Height  int
	Bar     progress.Model
	Console []string
}

type ConsoleWriter struct {
	channel chan<- ConsoleMsg
}

func NewProgressBar() (s3.ProgressCallback, error) {
	program := tea.NewProgram(
		&ProgressUI{
			Bar:     progress.New(progress.WithDefaultGradient()),
			Console: make([]string, 0),
		},
	)

	progressCh := make(chan ProgressMsg, 10)
	consoleCh := make(chan ConsoleMsg, 10)

	callback := func(current, total int) {
		progressCh <- ProgressMsg{current: current, total: total}
	}

	output := ConsoleWriter{channel: consoleCh}
	revert := logger.Temporary(output, logger.WithLevel(slog.LevelDebug))

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

	go func() {
		for msg := range consoleCh {
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
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			fmt.Println("SIGINT caught, upload canceled...")
			os.Exit(1)
		}
		return ui, nil

	case tea.WindowSizeMsg:
		ui.Width, ui.Height = msg.Width, msg.Height
		ui.Bar.Width = min(msg.Width, MaxWidth)
		return ui, nil

	case ProgressMsg:
		if ui.Bar.Percent() >= 1.0 {
			return ui, tea.Quit
		}

		pct := float64(msg.current) / float64(msg.total)
		cmd := ui.Bar.SetPercent(pct)
		return ui, cmd

	case ConsoleMsg:
		ui.Console = append(ui.Console, string(msg))
		return ui, nil

	case progress.FrameMsg:
		progressModel, cmd := ui.Bar.Update(msg)
		ui.Bar = progressModel.(progress.Model)
		return ui, cmd

	default:
		return ui, nil
	}
}

func (ui ProgressUI) View() string {
	wrappedLogs := lipgloss.NewStyle().
		Width(ui.Width).
		Render(strings.Join(ui.Console, "\n"))

	barStyle := lipgloss.NewStyle().
		Margin(1, 0, 2, 0).
		Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wrappedLogs,
		barStyle(ui.Bar.View()),
	)
}

func (w ConsoleWriter) Write(buf []byte) (n int, err error) {
	msg := strings.TrimRight(string(buf), "\n")
	w.channel <- ConsoleMsg(msg)
	return len(buf), nil
}
