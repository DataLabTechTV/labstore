package progressbar

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/IllumiKnowLabs/labstore/client/types"
	"github.com/IllumiKnowLabs/labstore/server/logger"
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

type Model struct {
	Ctx            context.Context
	Bar            progress.Model
	MaxConsoleSize int
	Debug          bool

	Progress chan types.Progress
	Message  chan string
	done     chan struct{}

	cancel  context.CancelFunc
	program *tea.Program
	width   int
	height  int
	console []string
}

type consoleWriter struct {
	ctx context.Context
	ch  chan<- string
}

func New(ctx context.Context, debug bool) (*Model, error) {
	ctx, cancel := context.WithCancel(ctx)

	m := &Model{
		Ctx:            ctx,
		Bar:            progress.New(progress.WithDefaultGradient()),
		MaxConsoleSize: defaultMaxConsoleSize,
		Debug:          debug,

		Progress: make(chan types.Progress, 10),
		Message:  make(chan string, 10),

		cancel:  cancel,
		done:    make(chan struct{}),
		console: make([]string, 0),
	}

	m.program = tea.NewProgram(m)

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return m.Bar.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressMsg:
		pct := float64(msg.current) / float64(msg.total)
		cmd := m.Bar.SetPercent(pct)
		return m, cmd

	case consoleMsg:
		m.console = append(m.console, string(msg))
		if len(m.console) > m.MaxConsoleSize {
			m.console = m.console[len(m.console)-m.MaxConsoleSize:]
		}
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			fmt.Println("SIGINT caught, quitting...")
			return m, tea.Quit
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.Bar.Update(msg)
		m.Bar = progressModel.(progress.Model)
		if !m.Bar.IsAnimating() && m.Bar.Percent() >= 1.0 {
			return m, tea.Quit
		}
		return m, cmd

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.Bar.Width = min(msg.Width, 80)
		return m, nil

	default:
		return m, nil
	}
}

func (m *Model) View() string {
	rows := []string{}

	if len(m.console) > 0 {
		wrappedLogs := lipgloss.NewStyle().
			Width(m.width).
			MarginBottom(1).
			Render(strings.Join(m.console, "\n"))
		rows = append(rows, wrappedLogs)
	}

	barStyle := lipgloss.NewStyle().
		MarginBottom(2).
		Render

	rows = append(rows, barStyle(m.Bar.View()))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)
}

func (m *Model) Run() {
	defer close(m.done)

	var wg sync.WaitGroup
	wg.Add(2)

	output := &consoleWriter{ctx: m.Ctx, ch: m.Message}
	revert := logger.Swap(output, logger.WithDebugFlag(m.Debug))

	go func() {
		defer wg.Done()
		for {
			select {
			case <-m.Ctx.Done():
				for {
					// Consume remaining
					select {
					case msg, ok := <-m.Progress:
						if !ok {
							return
						}
						m.program.Send(progressMsg{current: msg.Current, total: msg.Total})
					default:
						return
					}
				}

			case msg, ok := <-m.Progress:
				if !ok {
					return
				}
				m.program.Send(progressMsg{current: msg.Current, total: msg.Total})
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-m.Ctx.Done():
				return
			case msg, ok := <-m.Message:
				if !ok {
					return
				}
				m.program.Send(consoleMsg(msg))
			}
		}
	}()

	_, err := m.program.Run()
	m.cancel()
	wg.Wait()
	revert()

	if err != nil {
		slog.Error("progress bar", "err", err)
		return
	}

	slog.Debug("progress bar done")
}

func (m *Model) Close() {
	<-m.done
	close(m.Message)
}

func (w *consoleWriter) Write(buf []byte) (n int, err error) {
	msg := strings.TrimRight(string(buf), "\n")

	select {
	case <-w.ctx.Done():
		return 0, io.ErrClosedPipe
	case w.ch <- msg:
		return len(buf), nil
	}
}
