package tui

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/logger"
	"github.com/IllumiKnowLabs/labstore/client/types"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultMaxConsoleSize = 10

type progressMsg struct {
	current int
	total   int
}

type consoleMsg string

type ProgressBarModel struct {
	Bar            progress.Model
	MaxConsoleSize int

	Progress chan types.Progress
	Message  chan string
	Done     chan struct{}

	program *tea.Program
	width   int
	height  int
	console []string
}

type ConsoleWriter struct {
	channel chan<- string
}

func NewProgressBarModel() (*ProgressBarModel, error) {
	m := &ProgressBarModel{
		Bar:            progress.New(progress.WithDefaultGradient()),
		MaxConsoleSize: defaultMaxConsoleSize,

		Progress: make(chan types.Progress, 10),
		Message:  make(chan string, 10),
		Done:     make(chan struct{}),

		console: make([]string, 0),
	}

	m.program = tea.NewProgram(m)

	return m, nil
}

func (m *ProgressBarModel) Init() tea.Cmd {
	return m.Bar.Init()
}

func (m *ProgressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressMsg:
		pct := float64(msg.current) / float64(msg.total)
		cmd := m.Bar.SetPercent(pct)

		if pct >= 1.0 {
			return m, tea.Sequence(cmd, tea.Quit)
		}

		return m, cmd

	case consoleMsg:
		m.console = append(m.console, string(msg))
		if len(m.console) > m.MaxConsoleSize {
			m.console = m.console[len(m.console)-m.MaxConsoleSize:]
		}
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			fmt.Println("SIGINT caught, upload canceled...")
			return m, tea.Quit
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.Bar.Update(msg)
		m.Bar = progressModel.(progress.Model)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.Bar.Width = min(msg.Width, MaxWidth)
		return m, nil

	default:
		return m, nil
	}
}

func (m *ProgressBarModel) View() string {
	wrappedLogs := lipgloss.NewStyle().
		Width(m.width).
		Render(strings.Join(m.console, "\n"))

	barStyle := lipgloss.NewStyle().
		Margin(1, 0, 2, 0).
		Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wrappedLogs,
		barStyle(m.Bar.View()),
	)
}

func (m *ProgressBarModel) Run() {
	defer close(m.Done)

	var wg sync.WaitGroup
	wg.Add(2)

	output := ConsoleWriter{channel: m.Message}
	revert := logger.Temporary(output, logger.WithLevel(slog.LevelDebug))

	go func() {
		defer wg.Done()
		for msg := range m.Progress {
			m.program.Send(progressMsg{
				current: msg.Current,
				total:   msg.Total,
			})
		}
	}()

	go func() {
		defer wg.Done()
		for msg := range m.Message {
			m.program.Send(consoleMsg(msg))
		}
	}()

	_, err := m.program.Run()
	wg.Wait()
	revert()

	if err != nil {
		slog.Error("progress bar", "err", err)
		return
	}

	slog.Debug("progress bar done")
}

func (m *ProgressBarModel) Close() {
	close(m.Progress)
	close(m.Message)
	<-m.Done
}

func (w ConsoleWriter) Write(buf []byte) (n int, err error) {
	msg := strings.TrimRight(string(buf), "\n")
	w.channel <- msg
	return len(buf), nil
}
