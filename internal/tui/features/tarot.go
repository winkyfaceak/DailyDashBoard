package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tarotTextWidth = 70
	tarotImageCols = tarotTextWidth
	tarotImageRows = 16
)

// DailyTarot draws a three-card tarot spread.
func DailyTarot() *Feature {
	return &Feature{
		MenuLabel: "daily tarot",
		Loading:   "Drawing a three-card tarot spread...",
		ErrLabel:  "Tarot",
		BoxWidth:  84,
		Fetch: func() (Session, error) {
			reading, err := daily.RandomTarotReading()
			if err != nil {
				return nil, err
			}

			return &TarotSession{
				Reading:    reading,
				imageState: newImageState(renderTarotImage(reading)),
			}, nil
		},
	}
}

// TarotSession renders a tarot reading with card images.
type TarotSession struct {
	imageState
	Reading daily.TarotReading
}

func (s *TarotSession) View() string {
	return FormatTarotReading(s.Reading, s.imageText("Rendering tarot cards..."))
}

func (s *TarotSession) Links() []LinkAction {
	links := TarotImageLinks(s.Reading)
	if len(links) == 0 {
		return nil
	}

	return []LinkAction{{
		Key:     "o",
		URL:     links[0],
		Help:    "o to open first tarot card image",
		Message: "Opening first tarot card image...",
	}}
}

func renderTarotImage(reading daily.TarotReading) terminalimage.Image {
	var images []terminalimage.URLImage
	for _, card := range reading.Cards {
		if card.Card.ImageURL == "" {
			continue
		}

		images = append(images, terminalimage.URLImage{
			URL:        card.Card.ImageURL,
			UpsideDown: card.Orientation == daily.TarotOrientationReversed,
		})
	}

	if len(images) == 0 {
		return terminalimage.Image{Text: subtleTextStyle.Render("Tarot card images unavailable.")}
	}

	renderedImage, err := terminalimage.RenderGridFromURLs(images, tarotImageCols, tarotImageRows)
	if err != nil {
		return terminalimage.Image{Text: subtleTextStyle.Render(fmt.Sprintf("Tarot card images unavailable: %v", err))}
	}

	return renderedImage
}

func FormatTarotReading(reading daily.TarotReading, renderedImage string) string {
	var builder strings.Builder

	header := "Daily Tarot"
	if reading.SpreadName != "" {
		header += " -> " + reading.SpreadName
	}

	builder.WriteString(theme.Header(tarotTextWidth).Render(header))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n\n")
	}

	if len(reading.Cards) == 0 {
		builder.WriteString(formatSingleTarotCard(reading))
		return builder.String()
	}

	for index, card := range reading.Cards {
		if index > 0 {
			builder.WriteString("\n\n")
		}

		builder.WriteString(formatTarotSpreadCard(index, card))
	}

	if reading.Synthesis != "" {
		builder.WriteString("\n\n")
		builder.WriteString(labelStyle.Render("Synthesis"))
		builder.WriteString("\n")
		builder.WriteString(textutil.Wrap(reading.Synthesis, tarotTextWidth))
	}

	return builder.String()
}

func formatSingleTarotCard(reading daily.TarotReading) string {
	var builder strings.Builder

	builder.WriteString(strongTextStyle.Render(reading.Card.Name))

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

func formatTarotSpreadCard(index int, card daily.TarotSpreadCard) string {
	var builder strings.Builder

	position := card.Position
	if position == "" {
		position = fmt.Sprintf("Card %d", index+1)
	}

	builder.WriteString(labelStyle.Render(position))
	builder.WriteString(": ")
	builder.WriteString(strongTextStyle.Render(card.Card.Name))

	if card.Orientation != "" {
		builder.WriteString("  ")
		builder.WriteString(exampleStyle.Render(titleCase(card.Orientation)))
	}

	meta := formatTarotMeta(card.Card)
	if meta != "" {
		builder.WriteString("\n")
		builder.WriteString(subtleTextStyle.Render(meta))
	}

	if card.PositionFocus != "" {
		builder.WriteString("\n")
		builder.WriteString(exampleStyle.Render(textutil.Wrap(card.PositionFocus, tarotTextWidth)))
	}

	if card.Meaning != "" {
		builder.WriteString("\n")
		builder.WriteString(textutil.Wrap(card.Meaning, tarotTextWidth))
	}

	return builder.String()
}

func TarotImageLinks(reading daily.TarotReading) []string {
	var links []string
	for _, card := range reading.Cards {
		if card.Card.ImageURL != "" {
			links = append(links, card.Card.ImageURL)
		}
	}

	if len(links) == 0 && reading.Card.ImageURL != "" {
		links = append(links, reading.Card.ImageURL)
	}

	return links
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
