package tui

import (
	"context"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(ctx context.Context) {
	p := tea.NewProgram(NewTUIModel())
	if _, err := p.Run(); err != nil {
		slog.Error("tui model", "err", err)
		os.Exit(1)
	}
}
