//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly || solaris

package audio

import (
	"os"
	"syscall"
)

func pauseProcess(process *os.Process) error {
	return process.Signal(syscall.SIGSTOP)
}

func resumeProcess(process *os.Process) error {
	return process.Signal(syscall.SIGCONT)
}
