package tui

import tea "charm.land/bubbletea/v2"

type screen int

const (
	mainMenu screen = iota
	testMenu
)

type Model struct {
	screen           screen
	menus            map[screen]menu
	cursor           int
	message          string
	dictionaryResult string
	audioURL         string
}

func InitialModel() Model {
	return Model{
		screen: mainMenu,
		menus: map[screen]menu{
			mainMenu: {
				title: "Main Menu",
				items: []string{
					"test",
					"exit",
				},
			},
			testMenu: {
				title: "Test Menu",
				items: []string{
					"random word",
					"option 2",
					"back",
				},
			},
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
