package app

import (
	"DailyDashBoard/internal/tui/theme"

	"charm.land/lipgloss/v2"
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 3)

	spinnerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary)

	titleStyle = theme.Title

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

// contentBoxStyle is the bordered box a feature result renders in.
func contentBoxStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(1, 2).
		MarginTop(1).
		Width(width)
}
