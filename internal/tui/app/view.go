package app

import (
	"DailyDashBoard/internal/tui/features"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// minPageLines is the minimum content height kept visible on tiny terminals.
const minPageLines = 3

func (m Model) View() tea.View {
	var content string
	if m.screen == featurePage {
		// pageView composites its own full-screen background and centering.
		content = m.pageView()
	} else {
		content = m.menuView()

		if m.width > 0 && m.height > 0 {
			content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}
	}

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}

// menuView renders the main or daily menu centered in a bordered box.
func (m Model) menuView() string {
	currentMenu := m.currentMenu()

	var menuItems []string

	for i, choice := range currentMenu.items {
		if m.cursor == i {
			menuItems = append(menuItems, selectedItemStyle.Render("> "+choice.label))
		} else {
			menuItems = append(menuItems, normalItemStyle.Render("  "+choice.label))
		}
	}

	parts := []string{
		titleStyle.Render(currentMenu.title),
		"",
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
		helpStyle.Render("↑/↓ or j/k to move • enter to select • q to quit"),
	}

	if m.message != "" {
		parts = append(parts, messageStyle.Render(m.message))
	}

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

// pageView renders the dedicated page for the active feature: a loading
// spinner while fetching, then the session content with its help line.
// Content taller than the terminal is cropped and can be scrolled.
func (m Model) pageView() string {
	if m.loading {
		return m.loadingView()
	}

	var parts []string
	scrollable := false

	if m.session != nil && m.feature != nil {
		if content := m.session.View(); content != "" {
			box := contentBoxStyle(m.feature.BoxWidth, m.feature.Accent).Render(content)
			box, scrollable = m.cropToHeight(box)
			parts = append(parts, box)
		}
	}

	parts = append(parts, helpStyle.Render(m.pageHelpText(scrollable)))

	if m.message != "" {
		parts = append(parts, messageStyle.Render(m.message))
	}

	return m.compositePage(lipgloss.JoinVertical(lipgloss.Center, parts...))
}

// compositePage centers content within the terminal and, when a feature is
// active, layers it over that feature's scattered background. Falls back to
// plain content when the terminal size isn't known yet.
func (m Model) compositePage(content string) string {
	if m.width <= 0 || m.height <= 0 || m.feature == nil {
		return content
	}

	return composeBackground(m.width, m.height, content, m.feature.BackgroundWords, m.feature.Accent)
}

// pageFooterLines counts the lines rendered below the content box.
func (m Model) pageFooterLines() int {
	footer := 1 // help line
	if m.message != "" {
		footer++
	}

	return footer
}

// cropToHeight limits box to the lines that fit above the footer, applying
// the current scroll offset. It reports whether the box was cropped.
func (m Model) cropToHeight(box string) (string, bool) {
	if m.height <= 0 {
		return box, false
	}

	lines := strings.Split(box, "\n")

	available := m.height - m.pageFooterLines()
	if available < minPageLines {
		available = minPageLines
	}

	overflow := len(lines) - available
	if overflow <= 0 {
		return box, false
	}

	scroll := m.scroll
	if scroll > overflow {
		scroll = overflow
	}

	if scroll < 0 {
		scroll = 0
	}

	return strings.Join(lines[scroll:scroll+available], "\n"), true
}

// pageMaxScroll is the largest useful scroll offset for the current page.
func (m Model) pageMaxScroll() int {
	if m.loading || m.session == nil || m.feature == nil || m.height <= 0 {
		return 0
	}

	content := m.session.View()
	if content == "" {
		return 0
	}

	box := contentBoxStyle(m.feature.BoxWidth, m.feature.Accent).Render(content)
	lines := strings.Count(box, "\n") + 1

	available := m.height - m.pageFooterLines()
	if available < minPageLines {
		available = minPageLines
	}

	if overflow := lines - available; overflow > 0 {
		return overflow
	}

	return 0
}

func (m Model) loadingView() string {
	loadingText := "Loading..."
	spinner := spinnerStyle

	if m.feature != nil {
		loadingText = m.feature.Loading
		spinner = spinner.Foreground(m.feature.Accent)
	}

	frame := spinnerFrames[m.spinner%len(spinnerFrames)]

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		spinner.Render(frame)+" "+loadingText,
		"",
		helpStyle.Render("esc to go back • q to quit"),
	)

	return m.compositePage(content)
}

// pageHelpText builds the help line for a feature page: interactive sessions
// supply their own, otherwise navigation, audio, and link hints are combined.
func (m Model) pageHelpText(scrollable bool) string {
	if session, ok := m.session.(features.InteractiveSession); ok {
		return session.Help()
	}

	helpText := "esc to go back"

	if scrollable {
		helpText += " • ↑/↓ or j/k to scroll"
	}

	if m.audio.URL != "" {
		audioAction := "play"
		if m.audio.Starting {
			audioAction = "start"
		}

		if m.audio.Player != nil {
			audioAction = "pause"
			if m.audio.Paused {
				audioAction = "resume"
			}
		}

		helpText += " • p to " + audioAction + " " + m.currentAudioLabel()
	}

	for _, link := range m.sessionLinks() {
		helpText += " • " + link.Help
	}

	if len(m.storyLinks()) > 0 {
		helpText += " • 1-9 to open story"
	}

	return helpText + " • q to quit"
}
