//go:build windows

package audio

import (
	"fmt"
	"os"
)

func pauseProcess(process *os.Process) error {
	return fmt.Errorf("pause is not supported on this platform")
}

func resumeProcess(process *os.Process) error {
	return fmt.Errorf("resume is not supported on this platform")
}
