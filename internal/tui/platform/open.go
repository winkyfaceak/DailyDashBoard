package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func OpenURL(targetURL string) error {
	targetURL = strings.TrimSpace(targetURL)
	if targetURL == "" {
		return fmt.Errorf("no URL is available")
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", targetURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	default:
		if _, err := exec.LookPath("xdg-open"); err != nil {
			return fmt.Errorf("no URL opener was found")
		}

		cmd = exec.Command("xdg-open", targetURL)
	}

	return cmd.Start()
}
