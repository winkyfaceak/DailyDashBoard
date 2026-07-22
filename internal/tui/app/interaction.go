package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"
	"DailyDashBoard/internal/tui/platform"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type externalOpenMsg struct {
	Err error
}

func openExternalURLCmd(targetURL string) tea.Cmd {
	return func() tea.Msg {
		return externalOpenMsg{
			Err: platform.OpenURL(targetURL),
		}
	}
}

func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == "ctrl+c" {
		m.stopAudio()
		return m, tea.Quit
	}

	if session, ok := m.session.(features.InteractiveSession); ok {
		return m.handleSessionKey(session, msg)
	}

	if index, ok := features.NewsLinkIndex(key); ok {
		if stories := m.storyLinks(); len(stories) > 0 {
			return m.openStory(stories, index)
		}
	}

	switch key {
	case "q":
		m.stopAudio()
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		items := m.currentMenu().items
		if m.cursor < len(items)-1 {
			m.cursor++
		}

	case "p":
		return m.toggleAudio()

	case "enter":
		selected, ok := m.selectedMenuItem()
		if !ok {
			return m, nil
		}

		return m.handleSelection(selected)

	default:
		if link, ok := m.linkForKey(key); ok {
			m.message = link.Message
			return m, openExternalURLCmd(link.URL)
		}
	}

	return m, nil
}

// handleSessionKey routes a key press to the active interactive session and
// applies whatever the session asked for in response.
func (m Model) handleSessionKey(session features.InteractiveSession, msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	response := session.HandleKey(msg)

	if response.ToggleAudio {
		return m.toggleAudio()
	}

	if response.Clear {
		m.stopAudio()
		m = m.clearResultContent()
		m.message = response.Message

		return m, nil
	}

	if response.Handled {
		m.message = response.Message
	}

	return m, nil
}

func (m Model) openStory(stories []string, index int) (tea.Model, tea.Cmd) {
	if index >= len(stories) || stories[index] == "" {
		m.message = fmt.Sprintf("Story %d does not include a link.", index+1)
		return m, nil
	}

	m.message = fmt.Sprintf("Opening story %d...", index+1)
	return m, openExternalURLCmd(stories[index])
}

func (m Model) toggleAudio() (tea.Model, tea.Cmd) {
	if m.audio.URL == "" {
		return m, nil
	}

	audioLabel := m.currentAudioLabel()

	if m.audio.Player != nil {
		if m.audio.Paused {
			if err := m.audio.Player.Resume(); err != nil {
				m.message = fmt.Sprintf("Audio error: %v", err)
				return m, nil
			}

			m.audio.Paused = false
			m.message = "Resumed " + audioLabel + "."
			return m, nil
		}

		if err := m.audio.Player.Pause(); err != nil {
			m.message = fmt.Sprintf("Audio error: %v", err)
			return m, nil
		}

		m.audio.Paused = true
		m.message = "Paused " + audioLabel + "."
		return m, nil
	}

	if m.audio.Starting {
		m.message = "Starting " + audioLabel + "..."
		return m, nil
	}

	m.message = "Playing " + audioLabel + "..."
	m.audio.Starting = true
	m.audio.RequestID++

	return m, audio.StartCmd(m.audio.URL, m.audio.AllowDownload, m.audio.RequestID)
}

func (m Model) handleSelection(selected menuItem) (tea.Model, tea.Cmd) {
	switch selected.action {
	case actionOpenDailyMenu:
		return m.goTo(dailyMenu), nil
	case actionQuit:
		m.stopAudio()
		return m, tea.Quit
	case actionBack:
		return m.goTo(mainMenu), nil
	case actionRunFeature:
		if selected.feature == nil {
			break
		}

		m = m.clearContent()
		m.message = selected.feature.Loading

		return m, selected.feature.FetchCmd()
	}

	m = m.clearContent()
	m.message = fmt.Sprintf("No action is configured for %q.", selected.label)

	return m, nil
}
