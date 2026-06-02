package features

import "charm.land/lipgloss/v2"

var (
	primaryColor = lipgloss.Color("#7D56F4")
	successColor = lipgloss.Color("#04B575")
	textColor    = lipgloss.Color("#FAFAFA")
	mutedColor   = lipgloss.Color("#A8A8A8")
)

var (
	wordStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	partOfSpeechStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(successColor)

	factHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 1).
			MaxWidth(factTextWidth)

	exampleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(mutedColor)

	artworkHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 1).
				MaxWidth(artworkTextWidth)

	newsHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 1).
			MaxWidth(newsTextWidth)

	newsHeadlineStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor)

	musicHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 1).
				MaxWidth(musicTextWidth)

	musicTrackStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor)

	recipeHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 1).
				MaxWidth(recipeTextWidth)

	recipeNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor)

	subtleTextStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

func SubtleText(value string) string {
	return subtleTextStyle.Render(value)
}
