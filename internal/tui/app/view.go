package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	currentMenu := m.currentMenu()
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

	if m.links.ArtworkImage != "" {
		helpText += " • o to open image"
	}

	if m.links.MusicImage != "" {
		helpText += " • o to open album art"
	}

	if m.links.MusicTrack != "" {
		helpText += " • s to open song"
	}

	if m.links.RecipeImage != "" {
		helpText += " • o to open food image"
	}

	if m.links.RecipeSource != "" {
		helpText += " • r to open recipe"
	}

	if m.links.RecipeVideo != "" {
		helpText += " • v to open video"
	}

	if len(m.links.NewsStories) > 0 {
		helpText += " • 1-9 to open story"
	}

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
		helpStyle.Render(helpText),
	}

	if m.message != "" {
		contentParts = append(contentParts, messageStyle.Render(m.message))
	}

	if m.content.Fact != "" {
		contentParts = append(contentParts, factBoxStyle.Render(m.content.Fact))
	}

	if m.content.Dictionary != "" {
		contentParts = append(contentParts, dictionaryBoxStyle.Render(m.content.Dictionary))
	}

	if m.content.Tarot != "" {
		contentParts = append(contentParts, tarotBoxStyle.Render(m.content.Tarot))
	}

	if m.content.Artwork.Text != "" {
		contentParts = append(contentParts, artworkBoxStyle.Render(m.content.Artwork.Text))
	}

	if m.content.News != "" {
		contentParts = append(contentParts, newsBoxStyle.Render(m.content.News))
	}

	if m.content.Music.Text != "" {
		contentParts = append(contentParts, musicBoxStyle.Render(m.content.Music.Text))
	}

	if m.content.Recipe.Text != "" {
		contentParts = append(contentParts, recipeBoxStyle.Render(m.content.Recipe.Text))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentParts...,
	)

	box := boxStyle.Render(content)

	return tea.NewView(appStyle.Render(box))
}
