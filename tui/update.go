package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

func (m model) handleSelection(selected string) (tea.Model, tea.Cmd) {
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
		case "back":
			m = m.goTo(mainMenu)

		default:
			m.message = fmt.Sprintf("You selected from test menu: %s", selected)
		}
	}

	return m, nil
}
