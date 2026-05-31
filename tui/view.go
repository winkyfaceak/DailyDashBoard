package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	currentMenu := m.currentMenu()
	helpText := "↑/↓ or j/k to move • enter to select • q to quit"

	if m.audioURL != "" {
		helpText += " • p to play audio"
	}

	if m.artworkImageURL != "" {
		helpText += " • o to open image"
	}

	var menuItems []string

	for i, choice := range currentMenu.items {
		if m.cursor == i {
			menuItems = append(menuItems, selectedItemStyle.Render("> "+choice))
		} else {
			menuItems = append(menuItems, normalItemStyle.Render("  "+choice))
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
		helpStyle.Render(helpText),
	}

	if m.message != "" {
		contentParts = append(contentParts, messageStyle.Render(m.message))
	}

	if m.dictionaryResult != "" {
		contentParts = append(contentParts, dictionaryBoxStyle.Render(m.dictionaryResult))
	}

	if m.artworkResult != "" {
		contentParts = append(contentParts, artworkBoxStyle.Render(m.artworkResult))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentParts...,
	)

	box := boxStyle.Render(content)

	return tea.NewView(appStyle.Render(box))
}
