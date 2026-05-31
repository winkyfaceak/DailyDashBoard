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

func fetchDictionaryCmd(language string, word string) tea.Cmd {
	return func() tea.Msg {
		entries, err := api.DictionarySet(language, word)

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
			return m, nil
		}

		if len(msg.Entries) == 0 {
			m.message = "No dictionary entries found."
			m.dictionaryResult = ""
			return m, nil
		}

		m.message = ""
		m.dictionaryResult = formatDictionaryEntry(msg.Entries[0])

		return m, nil

	case tea.KeyPressMsg:
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
		case "word":
			m.message = "Searching dictionary..."
			m.dictionaryResult = ""
			return m, fetchDictionaryCmd("en", "hello")

		case "back":
			m = m.goTo(mainMenu)

		default:
			m.message = fmt.Sprintf("You selected from test menu: %s", selected)
			m.dictionaryResult = ""
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
		builder.WriteString(audio)
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
