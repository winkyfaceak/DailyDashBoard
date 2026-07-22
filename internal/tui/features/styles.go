package features

import (
	"DailyDashBoard/internal/tui/theme"

	"charm.land/lipgloss/v2"
)

// Shared styles, all derived from the theme package so the whole app uses a
// single palette. Feature headers use theme.Header(width) directly.
var (
	wordStyle = theme.Title

	labelStyle = theme.Label

	partOfSpeechStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.Success)

	exampleStyle = theme.Example

	subtleTextStyle = theme.Subtle

	successTextStyle = theme.SuccessText

	strongTextStyle = theme.Emphasis
)

// SubtleText renders muted secondary text.
func SubtleText(value string) string {
	return subtleTextStyle.Render(value)
}
