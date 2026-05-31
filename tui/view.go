package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	currentMenu := m.currentMenu()

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
		helpStyle.Render("↑/↓ or j/k to move • enter to select • q to quit"),
	}

	if m.message != "" {
		contentParts = append(contentParts, messageStyle.Render(m.message))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentParts...,
	)

	box := boxStyle.Render(content)

	return tea.NewView(appStyle.Render(box))
}
