package tui

import (
	"DailyDashBoard/dailys/api"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type dictionaryResultMsg struct {
	Entries []api.DictionaryEntry
	Err     error
}

type externalOpenMsg struct {
	Err error
}

func fetchRandomDictionaryCmd(language string) tea.Cmd {
	return func() tea.Msg {
		entries, err := api.RandomDictionarySet(language)

		return dictionaryResultMsg{
			Entries: entries,
			Err:     err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case factResultMsg:
		m = m.clearResultContent()

		if msg.Err != nil {
			m.message = fmt.Sprintf("Fact error: %v", msg.Err)
			return m, nil
		}

		m.message = ""
		m.factResult = formatFact(msg.Fact)

		return m, nil

	case dictionaryResultMsg:
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
		m.audioURL = firstAudio(entry.Phonetics)
		m.audioLabel = "pronunciation"
		m.audioAllowDownload = true
		m.dictionaryResult = formatDictionaryEntry(entry)

		return m, nil

	case artworkResultMsg:
		m = m.clearResultContent()

		if msg.Err != nil {
			m.message = fmt.Sprintf("Artwork error: %v", msg.Err)
			return m, nil
		}

		m.message = ""
		m.artworkImageURL = preferredArtworkImageURL(msg.Artwork)

		readyResult := formatArtwork(msg.Artwork, msg.Image.Text)
		if msg.Image.Raw != "" {
			m.artworkResult = formatArtwork(msg.Artwork, subtleTextStyle.Render("Rendering high-quality image..."))
			m.pendingArtwork = readyResult

			return m, tea.Sequence(
				tea.Raw(msg.Image.Raw),
				func() tea.Msg { return artworkImageReadyMsg{} },
			)
		}

		m.artworkResult = readyResult

		return m, nil

	case newsResultMsg:
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
		m.newsResult = formatNewsItems(msg.Items)
		m.newsLinks = newsLinks(msg.Items)

		return m, nil

	case musicResultMsg:
		m = m.clearResultContent()

		if msg.Err != nil {
			m.message = fmt.Sprintf("Music error: %v", msg.Err)
			return m, nil
		}

		m.message = ""
		m.musicImageURL = msg.Pick.ArtworkURL
		m.musicTrackURL = msg.Pick.TrackURL
		m.audioURL = msg.Pick.PreviewURL
		m.audioLabel = "preview"
		m.audioAllowDownload = true

		readyResult := formatMusic(msg.Pick, msg.Image.Text)
		if msg.Image.Raw != "" {
			m.musicResult = formatMusic(msg.Pick, subtleTextStyle.Render("Rendering album art..."))
			m.pendingMusic = readyResult

			return m, tea.Sequence(
				tea.Raw(msg.Image.Raw),
				func() tea.Msg { return musicImageReadyMsg{} },
			)
		}

		m.musicResult = readyResult
		m.pendingMusic = ""

		return m, nil

	case recipeResultMsg:
		m = m.clearResultContent()

		if msg.Err != nil {
			m.message = fmt.Sprintf("Recipe error: %v", msg.Err)
			return m, nil
		}

		m.message = ""
		m.recipeImageURL = msg.Recipe.ImageURL
		m.recipeSourceURL = msg.Recipe.SourceURL
		m.recipeVideoURL = msg.Recipe.YouTubeURL

		readyResult := formatRecipe(msg.Recipe, msg.Image.Text)
		if msg.Image.Raw != "" {
			m.recipeResult = formatRecipe(msg.Recipe, subtleTextStyle.Render("Rendering food image..."))
			m.pendingRecipe = readyResult

			return m, tea.Sequence(
				tea.Raw(msg.Image.Raw),
				func() tea.Msg { return recipeImageReadyMsg{} },
			)
		}

		m.recipeResult = readyResult
		m.pendingRecipe = ""

		return m, nil

	case artworkImageReadyMsg:
		if m.pendingArtwork != "" {
			m.artworkResult = m.pendingArtwork
			m.pendingArtwork = ""
		}

		return m, nil

	case musicImageReadyMsg:
		if m.pendingMusic != "" {
			m.musicResult = m.pendingMusic
			m.pendingMusic = ""
		}

		return m, nil

	case recipeImageReadyMsg:
		if m.pendingRecipe != "" {
			m.recipeResult = m.pendingRecipe
			m.pendingRecipe = ""
		}

		return m, nil

	case audioStartedMsg:
		if msg.RequestID != m.audioRequestID {
			if msg.Player != nil {
				msg.Player.Stop()
				return m, waitAudioCmd(msg.Player)
			}

			return m, nil
		}

		m.audioStarting = false

		if msg.Err != nil {
			m.message = fmt.Sprintf("Audio error: %v", msg.Err)
			return m, nil
		}

		if msg.Player == nil {
			m.message = "Audio error: no audio player was started"
			return m, nil
		}

		if m.audioPlayer != nil {
			m.audioPlayer.Stop()
		}

		m.audioPlayer = msg.Player
		m.audioPaused = false
		m.message = fmt.Sprintf("Playing %s.", m.currentAudioLabel())

		return m, waitAudioCmd(msg.Player)

	case audioFinishedMsg:
		if msg.Player != m.audioPlayer {
			return m, nil
		}

		m.audioPlayer = nil
		m.audioPaused = false
		m.audioStarting = false

		if msg.Err != nil {
			m.message = fmt.Sprintf("Audio error: %v", msg.Err)
			return m, nil
		}

		m.message = fmt.Sprintf("Finished %s.", m.currentAudioLabel())

		return m, nil

	case externalOpenMsg:
		if msg.Err != nil {
			m.message = fmt.Sprintf("Open link error: %v", msg.Err)
			return m, nil
		}

		m.message = "Opened link."

		return m, nil

	case tea.KeyPressMsg:
		if index, ok := newsLinkIndex(msg.String()); ok && len(m.newsLinks) > 0 {
			if index >= len(m.newsLinks) || m.newsLinks[index] == "" {
				m.message = fmt.Sprintf("Story %d does not include a link.", index+1)
				return m, nil
			}

			m.message = fmt.Sprintf("Opening story %d...", index+1)
			return m, openExternalURLCmd(m.newsLinks[index])
		}

		switch msg.String() {
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
			if m.audioURL != "" {
				audioLabel := m.currentAudioLabel()

				if m.audioPlayer != nil {
					if m.audioPaused {
						if err := m.audioPlayer.Resume(); err != nil {
							m.message = fmt.Sprintf("Audio error: %v", err)
							return m, nil
						}

						m.audioPaused = false
						m.message = "Resumed " + audioLabel + "."
						return m, nil
					}

					if err := m.audioPlayer.Pause(); err != nil {
						m.message = fmt.Sprintf("Audio error: %v", err)
						return m, nil
					}

					m.audioPaused = true
					m.message = "Paused " + audioLabel + "."
					return m, nil
				}

				if m.audioStarting {
					m.message = "Starting " + audioLabel + "..."
					return m, nil
				}

				m.message = "Playing " + audioLabel + "..."
				m.audioStarting = true
				m.audioRequestID++
				return m, startAudioCmd(m.audioURL, m.audioAllowDownload, m.audioRequestID)
			}

		case "o":
			if m.artworkImageURL != "" {
				m.message = "Opening full artwork image..."
				return m, openExternalURLCmd(m.artworkImageURL)
			}

			if m.musicImageURL != "" {
				m.message = "Opening album art..."
				return m, openExternalURLCmd(m.musicImageURL)
			}

			if m.recipeImageURL != "" {
				m.message = "Opening food image..."
				return m, openExternalURLCmd(m.recipeImageURL)
			}

		case "r":
			if m.recipeSourceURL != "" {
				m.message = "Opening recipe source..."
				return m, openExternalURLCmd(m.recipeSourceURL)
			}

		case "v":
			if m.recipeVideoURL != "" {
				m.message = "Opening cooking video..."
				return m, openExternalURLCmd(m.recipeVideoURL)
			}

		case "s":
			if m.musicTrackURL != "" {
				m.message = "Opening song..."
				return m, openExternalURLCmd(m.musicTrackURL)
			}

		case "enter":
			items := m.currentMenu().items
			selected := items[m.cursor]

			return m.handleSelection(selected)
		}
	}

	return m, nil
}

func (m Model) handleSelection(selected string) (tea.Model, tea.Cmd) {
	switch m.screen {
	case mainMenu:
		switch selected {
		case "test":
			m = m.goTo(testMenu)

		case "exit":
			return m, tea.Quit
		}

	case testMenu:
		switch selected {
		case "random word":
			m = m.clearContent()
			m.message = "Searching for a random word..."
			return m, fetchRandomDictionaryCmd("en")

		case "random fact":
			m = m.clearContent()
			m.message = "Fetching a random fact..."
			return m, fetchFactCmd()

		case "random artwork":
			m = m.clearContent()
			m.message = "Fetching artwork..."
			return m, fetchArtworkCmd()

		case "daily news":
			m = m.clearContent()
			m.message = "Fetching daily news..."
			return m, fetchNewsCmd()

		case "random music":
			m = m.clearContent()
			m.message = "Finding a genre and song..."
			return m, fetchMusicCmd()

		case "random recipe":
			m = m.clearContent()
			m.message = "Finding a random recipe..."
			return m, fetchRecipeCmd()

		case "back":
			m = m.goTo(mainMenu)

		default:
			m = m.clearContent()
			m.message = fmt.Sprintf("You selected from test menu: %s", selected)
		}
	}

	return m, nil
}

func (m Model) clearContent() Model {
	m.stopAudio()
	m.audioRequestID++

	m = m.clearResultContent()
	m.audioStarting = false

	return m
}

func (m Model) clearResultContent() Model {
	m.factResult = ""
	m.dictionaryResult = ""
	m.artworkResult = ""
	m.pendingArtwork = ""
	m.artworkImageURL = ""
	m.musicResult = ""
	m.pendingMusic = ""
	m.musicImageURL = ""
	m.musicTrackURL = ""
	m.newsResult = ""
	m.newsLinks = nil
	m.recipeResult = ""
	m.pendingRecipe = ""
	m.recipeImageURL = ""
	m.recipeSourceURL = ""
	m.recipeVideoURL = ""
	m.audioURL = ""
	m.audioLabel = ""
	m.audioAllowDownload = false

	return m
}

func (m *Model) stopAudio() {
	if m.audioPlayer != nil {
		m.audioPlayer.Stop()
		m.audioPlayer = nil
	}

	m.audioPaused = false
	m.audioStarting = false
}

func (m Model) currentAudioLabel() string {
	if m.audioLabel != "" {
		return m.audioLabel
	}

	return "audio"
}

func formatDictionaryEntry(entry api.DictionaryEntry) string {
	var builder strings.Builder

	builder.WriteString(wordStyle.Render(entry.Word))

	if entry.Phonetic != "" {
		builder.WriteString("  ")
		builder.WriteString(exampleStyle.Render(entry.Phonetic))
	}

	builder.WriteString("\n")

	if entry.Origin != "" {
		builder.WriteString("\n")
		builder.WriteString(labelStyle.Render("Origin: "))
		builder.WriteString(entry.Origin)
		builder.WriteString("\n")
	}

	audio := firstAudio(entry.Phonetics)
	if audio != "" {
		builder.WriteString(labelStyle.Render("Audio: "))
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(audio).Render("play pronunciation"))
		builder.WriteString("\n")
	}

	meaningCount := 0

	for _, meaning := range entry.Meanings {
		if len(meaning.Definitions) == 0 {
			continue
		}

		meaningCount++

		if meaningCount > 3 {
			break
		}

		builder.WriteString("\n")
		builder.WriteString(partOfSpeechStyle.Render(meaning.PartOfSpeech))
		builder.WriteString("\n")

		for i, definition := range meaning.Definitions {
			if i >= 2 {
				break
			}

			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, definition.Definition))

			if definition.Example != "" {
				builder.WriteString(exampleStyle.Render("     e.g. " + definition.Example))
				builder.WriteString("\n")
			}

			if len(definition.Synonyms) > 0 {
				builder.WriteString(labelStyle.Render("     Synonyms: "))
				builder.WriteString(strings.Join(limitStrings(definition.Synonyms, 5), ", "))
				builder.WriteString("\n")
			}

			if len(definition.Antonyms) > 0 {
				builder.WriteString(labelStyle.Render("     Antonyms: "))
				builder.WriteString(strings.Join(limitStrings(definition.Antonyms, 5), ", "))
				builder.WriteString("\n")
			}
		}
	}

	if meaningCount == 0 {
		builder.WriteString("\nNo meanings were available for this word.")
	}

	return builder.String()
}

func firstAudio(phonetics []api.Phonetic) string {
	for _, phonetic := range phonetics {
		if phonetic.Audio != "" {
			return phonetic.Audio
		}
	}

	return ""
}

func limitStrings(items []string, limit int) []string {
	if len(items) <= limit {
		return items
	}

	return items[:limit]
}
