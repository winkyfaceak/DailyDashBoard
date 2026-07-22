// Package theme holds the shared color palette and reusable lipgloss styles
// for the whole TUI, so every package renders with the same look.
package theme

import "charm.land/lipgloss/v2"

var (
	Primary = lipgloss.Color("#7D56F4")
	Success = lipgloss.Color("#04B575")
	Text    = lipgloss.Color("#FAFAFA")
	Muted   = lipgloss.Color("#A8A8A8")
	Help    = lipgloss.Color("#626262")
)

var (
	// Title is the bold inverted banner used for app and feature headers.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Text).
		Background(Primary).
		Padding(0, 1)

	// Label is a bold primary-colored field label ("Question", "Artist: ").
	Label = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary)

	// Emphasis is bold plain text for the main value on a line.
	Emphasis = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text)

	// Subtle is muted secondary text.
	Subtle = lipgloss.NewStyle().
		Foreground(Muted)

	// Example is muted italic text for examples and hints.
	Example = lipgloss.NewStyle().
		Italic(true).
		Foreground(Muted)

	// SuccessText marks correct answers and confirmations.
	SuccessText = lipgloss.NewStyle().
			Foreground(Success)
)

// Header returns the shared feature header style capped at maxWidth.
func Header(maxWidth int) lipgloss.Style {
	return Title.MaxWidth(maxWidth)
}
