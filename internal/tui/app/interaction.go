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

	if index, ok := features.NewsLinkIndex(key); ok && len(m.links.NewsStories) > 0 {
		return m.openNewsStory(index)
	}

	switch key {
	case "ctrl+c", "q":
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

	case "o":
		return m.openPrimaryImage()

	case "r":
		return m.openRecipeSource()

	case "v":
		return m.openRecipeVideo()

	case "s":
		return m.openMusicTrack()

	case "enter":
		selected, ok := m.selectedMenuItem()
		if !ok {
			return m, nil
		}

		return m.handleSelection(selected)
	}

	return m, nil
}

func (m Model) openNewsStory(index int) (tea.Model, tea.Cmd) {
	if index >= len(m.links.NewsStories) || m.links.NewsStories[index] == "" {
		m.message = fmt.Sprintf("Story %d does not include a link.", index+1)
		return m, nil
	}

	m.message = fmt.Sprintf("Opening story %d...", index+1)
	return m, openExternalURLCmd(m.links.NewsStories[index])
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

func (m Model) openPrimaryImage() (tea.Model, tea.Cmd) {
	if len(m.links.TarotImages) > 0 {
		m.message = "Opening first tarot card image..."
		return m, openExternalURLCmd(m.links.TarotImages[0])
	}

	if m.links.ArtworkImage != "" {
		m.message = "Opening full artwork image..."
		return m, openExternalURLCmd(m.links.ArtworkImage)
	}

	if m.links.MusicImage != "" {
		m.message = "Opening album art..."
		return m, openExternalURLCmd(m.links.MusicImage)
	}

	if m.links.RecipeImage != "" {
		m.message = "Opening food image..."
		return m, openExternalURLCmd(m.links.RecipeImage)
	}

	return m, nil
}

func (m Model) openRecipeSource() (tea.Model, tea.Cmd) {
	if m.links.RecipeSource == "" {
		return m, nil
	}

	m.message = "Opening recipe source..."
	return m, openExternalURLCmd(m.links.RecipeSource)
}

func (m Model) openRecipeVideo() (tea.Model, tea.Cmd) {
	if m.links.RecipeVideo == "" {
		return m, nil
	}

	m.message = "Opening cooking video..."
	return m, openExternalURLCmd(m.links.RecipeVideo)
}

func (m Model) openMusicTrack() (tea.Model, tea.Cmd) {
	if m.links.MusicTrack == "" {
		return m, nil
	}

	m.message = "Opening song..."
	return m, openExternalURLCmd(m.links.MusicTrack)
}

func (m Model) handleSelection(selected menuItem) (tea.Model, tea.Cmd) {
	switch selected.action {
	case actionOpenDailyMenu:
		return m.goTo(dailyMenu), nil
	case actionQuit:
		m.stopAudio()
		return m, tea.Quit
	case actionRandomWord:
		return m.startFeature("Searching for a random word...", features.FetchRandomDictionaryCmd("en"))
	case actionRandomFact:
		return m.startFeature("Fetching a random fact...", features.FetchFactCmd())
	case actionDailyTarot:
		return m.startFeature("Drawing a three-card tarot spread...", features.FetchTarotCmd())
	case actionRandomArtwork:
		return m.startFeature("Fetching artwork...", features.FetchArtworkCmd())
	case actionDailyNews:
		return m.startFeature("Fetching daily news...", features.FetchNewsCmd())
	case actionRandomMusic:
		return m.startFeature("Finding a genre and song...", features.FetchMusicCmd())
	case actionRandomRecipe:
		return m.startFeature("Finding a random recipe...", features.FetchRecipeCmd())
	case actionBack:
		return m.goTo(mainMenu), nil
	default:
		m = m.clearContent()
		m.message = fmt.Sprintf("No action is configured for %q.", selected.label)
		return m, nil
	}
}

func (m Model) startFeature(message string, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m = m.clearContent()
	m.message = message

	return m, cmd
}
