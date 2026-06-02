package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	currentMenu := m.currentMenu()
	helpText := "↑/↓ or j/k to move • enter to select • q to quit"

	if m.audioURL != "" {
		audioAction := "play"
		if m.audioStarting {
			audioAction = "start"
		}

		if m.audioPlayer != nil {
			audioAction = "pause"
			if m.audioPaused {
				audioAction = "resume"
			}
		}

		helpText += " • p to " + audioAction + " " + m.currentAudioLabel()
	}

	if m.artworkImageURL != "" {
		helpText += " • o to open image"
	}

	if m.musicImageURL != "" {
		helpText += " • o to open album art"
	}

	if m.musicTrackURL != "" {
		helpText += " • s to open song"
	}

	if m.recipeImageURL != "" {
		helpText += " • o to open food image"
	}

	if m.recipeSourceURL != "" {
		helpText += " • r to open recipe"
	}

	if m.recipeVideoURL != "" {
		helpText += " • v to open video"
	}

	if len(m.newsLinks) > 0 {
		helpText += " • 1-9 to open story"
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

	if m.factResult != "" {
		contentParts = append(contentParts, factBoxStyle.Render(m.factResult))
	}

	if m.dictionaryResult != "" {
		contentParts = append(contentParts, dictionaryBoxStyle.Render(m.dictionaryResult))
	}

	if m.artworkResult != "" {
		contentParts = append(contentParts, artworkBoxStyle.Render(m.artworkResult))
	}

	if m.newsResult != "" {
		contentParts = append(contentParts, newsBoxStyle.Render(m.newsResult))
	}

	if m.musicResult != "" {
		contentParts = append(contentParts, musicBoxStyle.Render(m.musicResult))
	}

	if m.recipeResult != "" {
		contentParts = append(contentParts, recipeBoxStyle.Render(m.recipeResult))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentParts...,
	)

	box := boxStyle.Render(content)

	return tea.NewView(appStyle.Render(box))
}
