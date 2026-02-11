package alert

import (
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarn
	AlertError
)

const DefaultWidth = 40

type Model struct {
	ID      string
	Title   string
	Message string
	Level   AlertLevel
	Width   int
}

var counter int = 0

func New(level AlertLevel, title, message string) (Model, tea.Cmd) {
	counter++
	id := fmt.Sprintf("alert-%d", counter)

	m := Model{
		ID:      id,
		Title:   title,
		Message: message,
		Level:   level,
		Width:   DefaultWidth,
	}

	cmd := tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return messages.AlertHideMsg{ID: m.ID}
	})

	return m, cmd
}
