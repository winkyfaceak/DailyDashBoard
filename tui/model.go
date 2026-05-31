package tui

import tea "charm.land/bubbletea/v2"

type screen int

const (
	mainMenu screen = iota
	testMenu
)

type model struct {
	screen  screen
	menus   map[screen]menu
	cursor  int
	message string
}

func InitialModel() model {
	return model{
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
					"time",
					"option 2",
					"back",
				},
			},
		},
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
