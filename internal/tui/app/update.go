package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// imageShownMsg fires after a session's raw terminal image has been flushed.
type imageShownMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case features.ResultMsg:
		return m.handleFeatureResult(msg)
	case imageShownMsg:
		return m.handleImageShown()
	case audio.StartedMsg:
		return m.handleAudioStarted(msg)
	case audio.FinishedMsg:
		return m.handleAudioFinished(msg)
	case externalOpenMsg:
		return m.handleExternalOpen(msg)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	default:
		return m, nil
	}
}

func (m Model) handleFeatureResult(msg features.ResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		var notice features.Notice
		if errors.As(msg.Err, &notice) {
			m.message = string(notice)
		} else {
			m.message = fmt.Sprintf("%s error: %v", msg.Feature.ErrLabel, msg.Err)
		}

		return m, nil
	}

	m.message = ""
	m.feature = msg.Feature
	m.session = msg.Session

	if session, ok := msg.Session.(features.AudioSession); ok {
		spec := session.Audio()
		m.audio.URL = spec.URL
		m.audio.Label = spec.Label
		m.audio.AllowDownload = spec.AllowDownload
	}

	if session, ok := msg.Session.(features.ImageSession); ok {
		if raw := session.RawImage(); raw != "" {
			return m, tea.Sequence(
				tea.Raw(raw),
				func() tea.Msg { return imageShownMsg{} },
			)
		}
	}

	return m, nil
}

func (m Model) handleImageShown() (tea.Model, tea.Cmd) {
	if session, ok := m.session.(features.ImageSession); ok {
		session.ImageShown()
	}

	return m, nil
}

func (m Model) handleAudioStarted(msg audio.StartedMsg) (tea.Model, tea.Cmd) {
	if msg.RequestID != m.audio.RequestID {
		if msg.Player != nil {
			msg.Player.Stop()
			return m, audio.WaitCmd(msg.Player)
		}

		return m, nil
	}

	m.audio.Starting = false

	if msg.Err != nil {
		m.message = fmt.Sprintf("Audio error: %v", msg.Err)
		return m, nil
	}

	if msg.Player == nil {
		m.message = "Audio error: no audio player was started"
		return m, nil
	}

	if m.audio.Player != nil {
		m.audio.Player.Stop()
	}

	m.audio.Player = msg.Player
	m.audio.Paused = false
	m.message = fmt.Sprintf("Playing %s.", m.currentAudioLabel())

	return m, audio.WaitCmd(msg.Player)
}

func (m Model) handleAudioFinished(msg audio.FinishedMsg) (tea.Model, tea.Cmd) {
	if msg.Player != m.audio.Player {
		return m, nil
	}

	m.audio.Player = nil
	m.audio.Paused = false
	m.audio.Starting = false

	if msg.Err != nil {
		m.message = fmt.Sprintf("Audio error: %v", msg.Err)
		return m, nil
	}

	m.message = fmt.Sprintf("Finished %s.", m.currentAudioLabel())

	return m, nil
}

func (m Model) handleExternalOpen(msg externalOpenMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.message = fmt.Sprintf("Open link error: %v", msg.Err)
		return m, nil
	}

	m.message = "Opened link."

	return m, nil
}
