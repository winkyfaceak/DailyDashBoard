package app

import (
	"DailyDashBoard/internal/tui/theme"
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	spinnerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary)

	titleStyle = theme.Title

	subtitleStyle = lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true).
			MarginBottom(1)

	normalItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedItemStyle = theme.Title

	helpStyle = lipgloss.NewStyle().
			Foreground(theme.Help).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(theme.Success).
			MarginTop(1)
)

// contentBoxStyle is the bordered box a feature result renders in, bordered
// in the feature's own accent color.
func contentBoxStyle(width int, accent color.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 2).
		MarginTop(1).
		Width(width)
}

// menuBoxStyle is the bordered box a menu screen renders in, bordered in
// accent (the app's purple by default, or the highlighted feature's own
// accent on the daily menu).
func menuBoxStyle(accent color.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 3)
}

// selectedFeatureItemStyle highlights a daily menu entry in its own
// feature's accent color.
func selectedFeatureItemStyle(accent color.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Text).
		Background(accent).
		Padding(0, 1)
}

// featureIconStyle tints an unselected daily menu entry's icon in its
// feature's accent color.
func featureIconStyle(accent color.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accent)
}
