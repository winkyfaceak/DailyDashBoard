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

type audioPlaybackMsg struct {
	Err error
}

type externalOpenMsg struct {
	Err error
}

func fetchDictionaryCmd(language string, word string) tea.Cmd {
	return func() tea.Msg {
		entries, err := api.DictionarySet(language, word)

		return dictionaryResultMsg{
			Entries: entries,
			Err:     err,
		}
	}
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

	case dictionaryResultMsg:
		if msg.Err != nil {
			m.message = fmt.Sprintf("Dictionary error: %v", msg.Err)
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
			return m, nil
		}

		if len(msg.Entries) == 0 {
			m.message = "No dictionary entries found."
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
			return m, nil
		}

		entry := msg.Entries[0]
		m.message = ""
		m.audioURL = firstAudio(entry.Phonetics)
		m.dictionaryResult = formatDictionaryEntry(entry)
		m.artworkResult = ""
		m.pendingArtwork = ""
		m.artworkImageURL = ""
		m.newsResult = ""
		m.newsLinks = nil

		return m, nil

	case artworkResultMsg:
		m.audioURL = ""
		m.newsResult = ""
		m.newsLinks = nil

		if msg.Err != nil {
			m.message = fmt.Sprintf("Artwork error: %v", msg.Err)
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			return m, nil
		}

		m.message = ""
		m.dictionaryResult = ""
		m.pendingArtwork = ""
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
		m.audioURL = ""
		m.dictionaryResult = ""
		m.artworkResult = ""
		m.pendingArtwork = ""
		m.artworkImageURL = ""
		m.newsResult = ""
		m.newsLinks = nil

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

	case artworkImageReadyMsg:
		if m.pendingArtwork != "" {
			m.artworkResult = m.pendingArtwork
			m.pendingArtwork = ""
		}

		return m, nil

	case audioPlaybackMsg:
		if msg.Err != nil {
			m.message = fmt.Sprintf("Audio error: %v", msg.Err)
			return m, nil
		}

		m.message = "Played pronunciation."

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
				m.message = "Playing pronunciation..."
				return m, playAudioCmd(m.audioURL)
			}

		case "o":
			if m.artworkImageURL != "" {
				m.message = "Opening full artwork image..."
				return m, openExternalURLCmd(m.artworkImageURL)
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
			m.message = "Searching for a random word..."
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
			return m, fetchRandomDictionaryCmd("en")

		case "random artwork":
			m.message = "Fetching artwork..."
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
			return m, fetchArtworkCmd()

		case "daily news":
			m.message = "Fetching daily news..."
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
			return m, fetchNewsCmd()

		case "back":
			m = m.goTo(mainMenu)

		default:
			m.message = fmt.Sprintf("You selected from test menu: %s", selected)
			m.dictionaryResult = ""
			m.artworkResult = ""
			m.pendingArtwork = ""
			m.artworkImageURL = ""
			m.newsResult = ""
			m.newsLinks = nil
			m.audioURL = ""
		}
	}

	return m, nil
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
