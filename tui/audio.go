package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

func playAudioCmd(audioURL string) tea.Cmd {
	return func() tea.Msg {
		return audioPlaybackMsg{
			Err: playAudio(audioURL),
		}
	}
}

func playAudio(audioURL string) error {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return fmt.Errorf("no pronunciation audio is available")
	}

	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("afplay"); err == nil {
			if err := downloadAndPlayWithAFPlay(audioURL); err == nil {
				return nil
			}
		}
	}

	for _, args := range audioPlayerCommands(audioURL) {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue
		}

		if err := exec.Command(args[0], args[1:]...).Run(); err == nil {
			return nil
		}
	}

	return openAudioURL(audioURL)
}

func downloadAndPlayWithAFPlay(audioURL string) error {
	req, err := http.NewRequest(http.MethodGet, audioURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("audio download failed: %s", resp.Status)
	}

	file, err := os.CreateTemp("", "dailydashboard-pronunciation-*.mp3")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return exec.Command("afplay", file.Name()).Run()
}

func audioPlayerCommands(audioURL string) [][]string {
	return [][]string{
		{"mpv", "--no-terminal", "--really-quiet", audioURL},
		{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", audioURL},
	}
}

func openAudioURL(audioURL string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", audioURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", audioURL)
	default:
		if _, err := exec.LookPath("xdg-open"); err != nil {
			return fmt.Errorf("no terminal audio player or URL opener was found")
		}

		cmd = exec.Command("xdg-open", audioURL)
	}

	return cmd.Start()
}
