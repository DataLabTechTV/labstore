//go:build !windows

package system

import "syscall"

func NewDaemonSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
