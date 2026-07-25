package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// RandomWord looks up a random word in the dictionary API.
func RandomWord() *Feature {
	return &Feature{
		MenuLabel:       "random word",
		Loading:         "Searching for a random word...",
		ErrLabel:        "Dictionary",
		BoxWidth:        70,
		Accent:          lipgloss.Color("#4EA8DE"),
		BackgroundWords: []string{"lexicon", "syntax", "syllable", "verbatim", "define", "§"},
		Icon:            "❝",
		Tagline:         "A word worth knowing, defined.",
		Fetch: func() (Session, error) {
			entries, err := daily.RandomDictionarySet("en")
			if err != nil {
				return nil, err
			}

			if len(entries) == 0 {
				return nil, Notice("No dictionary entries found.")
			}

			return &DictionarySession{Entry: entries[0]}, nil
		},
	}
}

// DictionarySession renders a dictionary entry with optional pronunciation.
type DictionarySession struct {
	Entry daily.DictionaryEntry
}

func (s *DictionarySession) View() string {
	return FormatDictionaryEntry(s.Entry)
}

func (s *DictionarySession) Audio() Audio {
	return Audio{
		URL:           FirstAudio(s.Entry.Phonetics),
		Label:         "pronunciation",
		AllowDownload: true,
	}
}

func FormatDictionaryEntry(entry daily.DictionaryEntry) string {
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

	audio := FirstAudio(entry.Phonetics)
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
				builder.WriteString(strings.Join(textutil.LimitStrings(definition.Synonyms, 5), ", "))
				builder.WriteString("\n")
			}

			if len(definition.Antonyms) > 0 {
				builder.WriteString(labelStyle.Render("     Antonyms: "))
				builder.WriteString(strings.Join(textutil.LimitStrings(definition.Antonyms, 5), ", "))
				builder.WriteString("\n")
			}
		}
	}

	if meaningCount == 0 {
		builder.WriteString("\nNo meanings were available for this word.")
	}

	return builder.String()
}

func FirstAudio(phonetics []daily.Phonetic) string {
	for _, phonetic := range phonetics {
		if phonetic.Audio != "" {
			return phonetic.Audio
		}
	}

	return ""
}
