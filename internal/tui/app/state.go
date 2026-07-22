package app

import "DailyDashBoard/internal/tui/features"

func (m Model) clearContent() Model {
	m.stopAudio()
	m.audio.RequestID++

	m = m.clearResultContent()
	m.audio.Starting = false

	return m
}

func (m Model) clearResultContent() Model {
	m.feature = nil
	m.session = nil
	m.audio.URL = ""
	m.audio.Label = ""
	m.audio.AllowDownload = false

	return m
}

func (m *Model) stopAudio() {
	if m.audio.Player != nil {
		m.audio.Player.Stop()
		m.audio.Player = nil
	}

	m.audio.Paused = false
	m.audio.Starting = false
}

func (m Model) currentAudioLabel() string {
	if m.audio.Label != "" {
		return m.audio.Label
	}

	return "audio"
}

// sessionLinks returns the active session's key-triggered links, if any.
func (m Model) sessionLinks() []features.LinkAction {
	if session, ok := m.session.(features.LinkSession); ok {
		return session.Links()
	}

	return nil
}

// storyLinks returns the active session's numbered story links, if any.
func (m Model) storyLinks() []string {
	if session, ok := m.session.(features.StorySession); ok {
		return session.StoryLinks()
	}

	return nil
}

func (m Model) linkForKey(key string) (features.LinkAction, bool) {
	for _, link := range m.sessionLinks() {
		if link.Key == key && link.URL != "" {
			return link, true
		}
	}

	return features.LinkAction{}, false
}
