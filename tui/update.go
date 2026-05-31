package tui

import (
	"DailyDashBoard/dailys/api"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type dictionaryResultMsg struct {
	Entries []api.DictionaryEntry
	Err     error
}

func fetchDictionaryCmd(language string, word string) tea.Cmd {
	return func() tea.Msg {
		entries, err := api.DictionarySet(language, word)

		return dictionaryResultMsg{
			Entries: entries,
			Err:     err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case dictionaryResultMsg:
		if msg.Err != nil {
			m.message = fmt.Sprintf("dictionary error: %s", msg.Err)
			return m, nil
		}

		if len(msg.Entries) == 0 {
			m.message = "no dictionary entries"
			return m, nil
		}

		entry := msg.Entries[0]

		for _, meaning := range entry.Meanings {
			for _, definition := range meaning.Definitions {
				if definition.Definition != "" {
					m.message = fmt.Sprintf(
						"%s (%s): %s",
						entry.Word,
						meaning.PartOfSpeech,
						definition.Definition,
					)
					return m, nil
				}
			}
		}
		m.message = fmt.Sprintf("Found %s, but no definition was available.", entry.Word)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			items := m.currentMenu().items

			if m.cursor < len(items)-1 {
				m.cursor++
			}

		case "enter":
			items := m.currentMenu().items
			selected := items[m.cursor]

			return m.handleSelection(selected)
		}
	}

	return m, nil
}

func (m Model) handleSelection(selected string) (tea.Model, tea.Cmd) {
	switch m.screen {
	case mainMenu:
		switch selected {
		case "test":
			m = m.goTo(testMenu)

		case "exit":
			return m, tea.Quit
		}

	case testMenu:
		switch selected {
		case "word":
			m.message = "Searching dictionary..."
			return m, fetchDictionaryCmd("en", "hello")

		case "back":
			m = m.goTo(mainMenu)

		default:
			m.message = fmt.Sprintf("You selected from test menu: %s", selected)
		}
	}

	return m, nil
}
