package tui

import (
	"context"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(ctx context.Context) {
	program := tea.NewProgram(
		NewTUIModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		slog.Error("tui model", "err", err)
		os.Exit(1)
	}
}
