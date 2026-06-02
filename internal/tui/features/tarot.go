package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"strings"
	"unicode"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
)

const tarotTextWidth = 70

type TarotResultMsg struct {
	Reading daily.TarotReading
	Err     error
}

func FetchTarotCmd() tea.Cmd {
	return func() tea.Msg {
		reading, err := daily.RandomTarotReading()

		return TarotResultMsg{
			Reading: reading,
			Err:     err,
		}
	}
}

func FormatTarotReading(reading daily.TarotReading) string {
	var builder strings.Builder

	builder.WriteString(tarotHeaderStyle.Render("Daily Tarot"))
	builder.WriteString("\n\n")
	builder.WriteString(tarotCardStyle.Render(reading.Card.Name))

	if reading.Orientation != "" {
		builder.WriteString("  ")
		builder.WriteString(exampleStyle.Render(titleCase(reading.Orientation)))
	}

	meta := formatTarotMeta(reading.Card)
	if meta != "" {
		builder.WriteString("\n")
		builder.WriteString(subtleTextStyle.Render(meta))
	}

	if reading.Meaning != "" {
		builder.WriteString("\n\n")
		builder.WriteString(labelStyle.Render("Reading"))
		builder.WriteString("\n")
		builder.WriteString(textutil.Wrap(reading.Meaning, tarotTextWidth))
	}

	if reading.Card.Description != "" {
		builder.WriteString("\n\n")
		builder.WriteString(labelStyle.Render("Card"))
		builder.WriteString("\n")
		builder.WriteString(exampleStyle.Render(textutil.Wrap(reading.Card.Description, tarotTextWidth)))
	}

	return builder.String()
}

func formatTarotMeta(card daily.TarotCard) string {
	var parts []string

	if card.Type != "" {
		parts = append(parts, titleCase(card.Type)+" arcana")
	}

	if card.Suit != "" {
		parts = append(parts, titleCase(card.Suit))
	}

	if card.Value != "" {
		parts = append(parts, titleCase(strings.ToLower(card.Value)))
	}

	return strings.Join(parts, " • ")
}

func titleCase(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError {
		return value
	}

	return string(unicode.ToUpper(first)) + strings.ToLower(value[size:])
}
