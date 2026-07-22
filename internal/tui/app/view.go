package app

import (
	"DailyDashBoard/internal/tui/features"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	currentMenu := m.currentMenu()

	var menuItems []string

	for i, choice := range currentMenu.items {
		if m.cursor == i {
			menuItems = append(menuItems, selectedItemStyle.Render("> "+choice.label))
		} else {
			menuItems = append(menuItems, normalItemStyle.Render("  "+choice.label))
		}
	}

	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		menuItems...,
	)

	contentParts := []string{
		titleStyle.Render(currentMenu.title),
		"",
		menu,
		helpStyle.Render(m.helpText()),
	}

	if m.message != "" {
		contentParts = append(contentParts, messageStyle.Render(m.message))
	}

	if m.session != nil && m.feature != nil {
		if content := m.session.View(); content != "" {
			contentParts = append(contentParts, contentBoxStyle(m.feature.BoxWidth).Render(content))
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentParts...,
	)

	box := boxStyle.Render(content)

	return tea.NewView(appStyle.Render(box))
}

// helpText builds the help line: interactive sessions supply their own,
// otherwise the menu help is extended with audio and link hints.
func (m Model) helpText() string {
	if session, ok := m.session.(features.InteractiveSession); ok {
		return session.Help()
	}

	helpText := "↑/↓ or j/k to move • enter to select • q to quit"

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

	return helpText
}
