package app

import "charm.land/lipgloss/v2"

var (
	primaryColor = lipgloss.Color("#7D56F4")
	successColor = lipgloss.Color("#04B575")
	textColor    = lipgloss.Color("#FAFAFA")
	mutedColor   = lipgloss.Color("#A8A8A8")
	helpColor    = lipgloss.Color("#626262")
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(100)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(helpColor).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(successColor).
			MarginTop(1)

	factBoxStyle        = contentBoxStyle(84)
	mathBoxStyle        = contentBoxStyle(78)
	triviaBoxStyle      = contentBoxStyle(84)
	dictionaryBoxStyle  = contentBoxStyle(70)
	tarotBoxStyle       = contentBoxStyle(84)
	artworkBoxStyle     = contentBoxStyle(84)
	newsBoxStyle        = contentBoxStyle(92)
	musicBoxStyle       = contentBoxStyle(84)
	musicTriviaBoxStyle = contentBoxStyle(84)
	recipeBoxStyle      = contentBoxStyle(84)
)

func contentBoxStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2).
		MarginTop(1).
		Width(width)
}
