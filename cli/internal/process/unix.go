//go:build !windows

package process

import (
	"os"
	"syscall"
)

func NewDaemonSysAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func Kill(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}
