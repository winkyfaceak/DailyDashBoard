package audio

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
)

type StartedMsg struct {
	Player    *Player
	RequestID int
	Err       error
}

type FinishedMsg struct {
	Player *Player
	Err    error
}

type Player struct {
	cmd      *exec.Cmd
	tempFile string

	mu      sync.Mutex
	stopped bool
	cleaned bool
}

func StartCmd(audioURL string, allowDownload bool, requestID int) tea.Cmd {
	return func() tea.Msg {
		player, err := startAudio(audioURL, allowDownload)

		return StartedMsg{
			Player:    player,
			RequestID: requestID,
			Err:       err,
		}
	}
}

func WaitCmd(player *Player) tea.Cmd {
	return func() tea.Msg {
		return FinishedMsg{
			Player: player,
			Err:    player.Wait(),
		}
	}
}

func startAudio(audioURL string, allowDownload bool) (*Player, error) {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return nil, fmt.Errorf("no audio is available")
	}

	var lastErr error

	if allowDownload && runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("afplay"); err == nil {
			tempFile, err := downloadAudioTempFile(audioURL)
			if err == nil {
				player, err := startAudioCommand("afplay", []string{tempFile}, tempFile)
				if err == nil {
					return player, nil
				}

				lastErr = err
				_ = os.Remove(tempFile)
			} else {
				lastErr = err
			}
		}
	}

	for _, args := range audioPlayerCommands(audioURL) {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue
		}

		player, err := startAudioCommand(args[0], args[1:], "")
		if err == nil {
			return player, nil
		}

		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("no audio player was found")
}

func downloadAudioTempFile(audioURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, audioURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("audio download failed: %s", resp.Status)
	}

	file, err := os.CreateTemp("", "dailydashboard-audio-*.m4a")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		_ = os.Remove(file.Name())
		return "", err
	}

	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

func startAudioCommand(name string, args []string, tempFile string) (*Player, error) {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Player{
		cmd:      cmd,
		tempFile: tempFile,
	}, nil
}

func (p *Player) Pause() error {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("no audio is playing")
	}

	return pauseProcess(p.cmd.Process)
}

func (p *Player) Resume() error {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("no audio is paused")
	}

	return resumeProcess(p.cmd.Process)
}

func (p *Player) Stop() {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return
	}

	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}

	p.stopped = true
	p.mu.Unlock()

	_ = p.cmd.Process.Kill()
}

func (p *Player) Wait() error {
	if p == nil || p.cmd == nil {
		return nil
	}

	err := p.cmd.Wait()
	p.cleanup()

	p.mu.Lock()
	stopped := p.stopped
	p.mu.Unlock()

	if stopped {
		return nil
	}

	return err
}

func (p *Player) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cleaned {
		return
	}

	p.cleaned = true
	if p.tempFile != "" {
		_ = os.Remove(p.tempFile)
	}
}

func audioPlayerCommands(audioURL string) [][]string {
	return [][]string{
		{"mpv", "--no-terminal", "--really-quiet", audioURL},
		{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", audioURL},
	}
}
