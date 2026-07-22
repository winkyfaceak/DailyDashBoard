package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"
	"errors"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

const spinnerInterval = 100 * time.Millisecond

// imageShownMsg fires after a session's raw terminal image has been flushed.
type imageShownMsg struct{}

// featureResultMsg wraps a feature fetch result with the fetch it belongs
// to, so results that arrive after the user backed out are discarded.
type featureResultMsg struct {
	fetchID int
	result  features.ResultMsg
}

type spinnerTickMsg struct {
	fetchID int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if maxScroll := m.pageMaxScroll(); m.scroll > maxScroll {
			m.scroll = maxScroll
		}

		return m, nil
	case featureResultMsg:
		return m.handleFeatureResult(msg)
	case spinnerTickMsg:
		return m.handleSpinnerTick(msg)
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

// fetchFeatureCmd runs the feature's fetch and tags the result with the
// current fetch ID.
func (m Model) fetchFeatureCmd(feature *features.Feature) tea.Cmd {
	fetchID := m.fetchID
	fetch := feature.FetchCmd()

	return func() tea.Msg {
		result, _ := fetch().(features.ResultMsg)

		return featureResultMsg{
			fetchID: fetchID,
			result:  result,
		}
	}
}

func (m Model) spinnerTickCmd() tea.Cmd {
	fetchID := m.fetchID

	return tea.Tick(spinnerInterval, func(time.Time) tea.Msg {
		return spinnerTickMsg{fetchID: fetchID}
	})
}

func (m Model) handleSpinnerTick(msg spinnerTickMsg) (tea.Model, tea.Cmd) {
	if !m.loading || msg.fetchID != m.fetchID {
		return m, nil
	}

	m.spinner++

	return m, m.spinnerTickCmd()
}

func (m Model) handleFeatureResult(msg featureResultMsg) (tea.Model, tea.Cmd) {
	if msg.fetchID != m.fetchID || m.screen != featurePage {
		return m, nil
	}

	m = m.clearResultContent()
	m.loading = false
	m.scroll = 0
	m.feature = msg.result.Feature

	if msg.result.Err != nil {
		var notice features.Notice
		if errors.As(msg.result.Err, &notice) {
			m.message = string(notice)
		} else {
			m.message = fmt.Sprintf("%s error: %v", msg.result.Feature.ErrLabel, msg.result.Err)
		}

		return m, nil
	}

	m.message = ""
	m.session = msg.result.Session

	if session, ok := msg.result.Session.(features.AudioSession); ok {
		spec := session.Audio()
		m.audio.URL = spec.URL
		m.audio.Label = spec.Label
		m.audio.AllowDownload = spec.AllowDownload
	}

	if session, ok := msg.result.Session.(features.ImageSession); ok {
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
