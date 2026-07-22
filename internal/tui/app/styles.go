package app

import (
	"DailyDashBoard/internal/tui/theme"

	"charm.land/lipgloss/v2"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			Width(100)

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
