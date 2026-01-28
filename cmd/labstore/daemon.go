package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/IllumiKnowLabs/labstore/cli/pkg/process"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/server/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	pidFilePath       = filepath.Join(os.TempDir(), "labstore.pid")
	daemonAnnotations = map[string]string{
		"mode": "daemon",
	}
)

func AddDaemonCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:         "start",
			Short:       "Start LabStore server",
			Run:         start,
			Annotations: daemonAnnotations,
		},
		&cobra.Command{
			Use:         "stop",
			Short:       "Stop LabStore server",
			Run:         stop,
			Annotations: daemonAnnotations,
		},
		&cobra.Command{
			Use:         "status",
			Short:       "Check if LabStore server is running",
			Run:         status,
			Annotations: daemonAnnotations,
		},
		&cobra.Command{
			Use:         "restart",
			Short:       "Restart LabStore server",
			Run:         restart,
			Annotations: daemonAnnotations,
		},
	)
}

func start(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		fmt.Println("error: unsupported in windows, use the serve command instead (ctrl+c to interrupt)")
		return
	}

	w, err := logger.NewDailyWriter(config.App.Log.Dir, strings.ToLower(constants.Name))
	if err != nil {
		fmt.Println("error: could not open log file")
		return
	}
	logger.InitWithWriter(w)

	slog.Info("labstore starting")

	if pid, running := readPID(); running {
		slog.Error("labstore already running", "pid", pid)
		fmt.Printf("LabStore already running (PID=%d)\n", pid)
		return
	}

	c := exec.Command(os.Args[0], "serve")
	c.Stdout = w
	c.Stderr = w
	c.SysProcAttr = process.NewDaemonSysAttr()

	if err := c.Start(); err != nil {
		slog.Error("failed to start", "err", err)
		fmt.Println("error: failed to start")
		return
	}

	pid := c.Process.Pid
	if err := os.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0o644); err != nil {
		slog.Error("failed to write pid file", "pid", pid, "path", pidFilePath, "err", err)
		fmt.Println("error: failed to write PID file")
		return
	}

	slog.Info("labstore started", "pid", pid)
	fmt.Printf("LabStore started (PID=%d)\n", pid)
}

func stop(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		fmt.Println("error: unsupported in windows, use the serve command instead (ctrl+c to interrupt)")
		return
	}

	w, err := logger.NewDailyWriter(config.App.Log.Dir, strings.ToLower(constants.Name))
	if err != nil {
		fmt.Println("error: could not open log file")
		return
	}
	logger.InitWithWriter(w)

	slog.Info("labstore stopping")

	pid, running := readPID()
	if !running {
		slog.Error("labstore not running")
		fmt.Println("LabStore not running")
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		slog.Error("process not found", "pid", pid, "err", err)
		fmt.Printf("error: process not found (PID=%d)\n", pid)
		return
	}

	if err := process.Kill(proc); err != nil {
		slog.Error("failed to stop process", "pid", pid, "err", err)
		fmt.Printf("error: failed to stop process (PID=%d)\n", pid)
		return
	}

	if err := os.Remove(pidFilePath); err != nil {
		slog.Warn("could not delete pid file", "path", pidFilePath, "err", err)
		fmt.Println("Warning: could not delete PID file:", pidFilePath)
	}
	slog.Info("labstore stopped", "pid", pid)
	fmt.Printf("LabStore stopped (PID=%d)\n", pid)
}

func restart(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		fmt.Println("error: unsupported in windows, use the serve command instead (ctrl+c to interrupt)")
		return
	}

	stop(cmd, args)
	start(cmd, args)
}

func status(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		fmt.Println("error: unsupported in windows, use the serve command instead (ctrl+c to interrupt)")
		return
	}

	w, err := logger.NewDailyWriter(config.App.Log.Dir, strings.ToLower(constants.Name))
	if err != nil {
		fmt.Println("error: could not open log file")
		return
	}
	logger.InitWithWriter(w)

	if pid, running := readPID(); running {
		fmt.Printf("LabStore running (PID=%d)\n", pid)
	} else {
		fmt.Println("LabStore not running")
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

	proc, err := os.FindProcess(pid)
	if err != nil {
		return 0, false
	}
	if err := proc.Kill(); err != nil {
		return 0, false
	}

	return pid, true
}
