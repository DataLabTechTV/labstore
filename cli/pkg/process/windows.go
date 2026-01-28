//go:build windows

package process

import (
	"os"
	"syscall"
)

func NewDaemonSysAttr() *syscall.SysProcAttr {
	return nil
}

func Kill(proc *os.Process) error {
	return proc.Signal(os.Interrupt)
}
