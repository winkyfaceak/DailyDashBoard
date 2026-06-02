package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type artworkImageReadyMsg struct{}
type musicImageReadyMsg struct{}
type recipeImageReadyMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case features.FactResultMsg:
		return m.handleFactResult(msg)
	case features.DictionaryResultMsg:
		return m.handleDictionaryResult(msg)
	case features.ArtworkResultMsg:
		return m.handleArtworkResult(msg)
	case features.NewsResultMsg:
		return m.handleNewsResult(msg)
	case features.MusicResultMsg:
		return m.handleMusicResult(msg)
	case features.RecipeResultMsg:
		return m.handleRecipeResult(msg)
	case artworkImageReadyMsg:
		return m.handleArtworkImageReady()
	case musicImageReadyMsg:
		return m.handleMusicImageReady()
	case recipeImageReadyMsg:
		return m.handleRecipeImageReady()
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

func (m Model) handleFactResult(msg features.FactResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("Fact error: %v", msg.Err)
		return m, nil
	}

	m.message = ""
	m.content.Fact = features.FormatFact(msg.Fact)

	return m, nil
}

func (m Model) handleDictionaryResult(msg features.DictionaryResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("Dictionary error: %v", msg.Err)
		return m, nil
	}

	if len(msg.Entries) == 0 {
		m.message = "No dictionary entries found."
		return m, nil
	}

	entry := msg.Entries[0]
	m.message = ""
	m.audio.URL = features.FirstAudio(entry.Phonetics)
	m.audio.Label = "pronunciation"
	m.audio.AllowDownload = true
	m.content.Dictionary = features.FormatDictionaryEntry(entry)

	return m, nil
}

func (m Model) handleArtworkResult(msg features.ArtworkResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("Artwork error: %v", msg.Err)
		return m, nil
	}

	m.message = ""
	m.links.ArtworkImage = features.PreferredArtworkImageURL(msg.Artwork)

	readyResult := features.FormatArtwork(msg.Artwork, msg.Image.Text)
	if msg.Image.Raw != "" {
		m.content.Artwork.Text = features.FormatArtwork(msg.Artwork, features.SubtleText("Rendering high-quality image..."))
		m.content.Artwork.Pending = readyResult

		return m, tea.Sequence(
			tea.Raw(msg.Image.Raw),
			func() tea.Msg { return artworkImageReadyMsg{} },
		)
	}

	m.content.Artwork.Text = readyResult

	return m, nil
}

func (m Model) handleNewsResult(msg features.NewsResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("News error: %v", msg.Err)
		return m, nil
	}

	if len(msg.Items) == 0 {
		m.message = "No news stories found."
		return m, nil
	}

	m.message = ""
	m.content.News = features.FormatNewsItems(msg.Items)
	m.links.NewsStories = features.NewsLinks(msg.Items)

	return m, nil
}

func (m Model) handleMusicResult(msg features.MusicResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("Music error: %v", msg.Err)
		return m, nil
	}

	m.message = ""
	m.links.MusicImage = msg.Pick.ArtworkURL
	m.links.MusicTrack = msg.Pick.TrackURL
	m.audio.URL = msg.Pick.PreviewURL
	m.audio.Label = "preview"
	m.audio.AllowDownload = true

	readyResult := features.FormatMusic(msg.Pick, msg.Image.Text)
	if msg.Image.Raw != "" {
		m.content.Music.Text = features.FormatMusic(msg.Pick, features.SubtleText("Rendering album art..."))
		m.content.Music.Pending = readyResult

		return m, tea.Sequence(
			tea.Raw(msg.Image.Raw),
			func() tea.Msg { return musicImageReadyMsg{} },
		)
	}

	m.content.Music.Text = readyResult

	return m, nil
}

func (m Model) handleRecipeResult(msg features.RecipeResultMsg) (tea.Model, tea.Cmd) {
	m = m.clearResultContent()

	if msg.Err != nil {
		m.message = fmt.Sprintf("Recipe error: %v", msg.Err)
		return m, nil
	}

	m.message = ""
	m.links.RecipeImage = msg.Recipe.ImageURL
	m.links.RecipeSource = msg.Recipe.SourceURL
	m.links.RecipeVideo = msg.Recipe.YouTubeURL

	readyResult := features.FormatRecipe(msg.Recipe, msg.Image.Text)
	if msg.Image.Raw != "" {
		m.content.Recipe.Text = features.FormatRecipe(msg.Recipe, features.SubtleText("Rendering food image..."))
		m.content.Recipe.Pending = readyResult

		return m, tea.Sequence(
			tea.Raw(msg.Image.Raw),
			func() tea.Msg { return recipeImageReadyMsg{} },
		)
	}

	m.content.Recipe.Text = readyResult

	return m, nil
}

func (m Model) handleArtworkImageReady() (tea.Model, tea.Cmd) {
	if m.content.Artwork.Pending != "" {
		m.content.Artwork.Text = m.content.Artwork.Pending
		m.content.Artwork.Pending = ""
	}

	return m, nil
}

func (m Model) handleMusicImageReady() (tea.Model, tea.Cmd) {
	if m.content.Music.Pending != "" {
		m.content.Music.Text = m.content.Music.Pending
		m.content.Music.Pending = ""
	}

	return m, nil
}

func (m Model) handleRecipeImageReady() (tea.Model, tea.Cmd) {
	if m.content.Recipe.Pending != "" {
		m.content.Recipe.Text = m.content.Recipe.Pending
		m.content.Recipe.Pending = ""
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
