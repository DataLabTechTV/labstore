package main

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/server/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	pidFilePath = filepath.Join(os.TempDir(), "labstore.pid")
)

func AddDaemonCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start LabStore server",
			Run:   start,
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop LabStore server",
			Run:   stop,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check if LabStore server is running",
			Run:   status,
		},
		&cobra.Command{
			Use:   "restart",
			Short: "Restart LabStore server",
			Run: func(cmd *cobra.Command, args []string) {
				stop(cmd, args)
				start(cmd, args)
			},
		},
	)
}

func start(cmd *cobra.Command, args []string) {
	w, err := logger.NewDailyWriter(config.App.Log.Dir, strings.ToLower(constants.Name))
	if err != nil {
		slog.Error("could not open log file", "err", err)
		return
	}
	logger.InitWithWriter(w)

	if pid, running := readPID(); running {
		slog.Error("labstore already running", "pid", pid)
		return
	}

	c := exec.Command(os.Args[0], "serve")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := c.Start(); err != nil {
		slog.Error("failed to start", "err", err)
		return
	}

	pid := c.Process.Pid
	if err := os.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0o644); err != nil {
		slog.Error("failed to write pid file", "pid", pid, "err", err)
		return
	}

	slog.Info("labstore started", "pid", pid)
}

func stop(cmd *cobra.Command, args []string) {
	pid, running := readPID()
	if !running {
		slog.Error("labstore not running")
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		slog.Error("process not found", "pid", pid, "err", err)
		return
	}

	if err := proc.Kill(); err != nil {
		slog.Error("failed to stop process", "pid", pid, "err", err)
		return
	}

	if err := os.Remove(pidFilePath); err != nil {
		slog.Warn("could not delete pid file", "path", pidFilePath, "err", err)
	}
	slog.Info("labstore stopped", "pid", pid)
}

func status(cmd *cobra.Command, args []string) {
	pid, running := readPID()
	if running {
		slog.Info("labstore running", "pid", pid)
	} else {
		slog.Info("labstore not running")
	}
}

func readPID() (int, bool) {
	data, err := os.ReadFile(pidFilePath)
	if err != nil {
		return 0, false
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, false
	}

	if err := syscall.Kill(pid, 0); err != nil {
		return 0, false
	}

	return pid, true
}
