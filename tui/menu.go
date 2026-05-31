package tui

type menu struct {
	title string
	items []string
}

func (m Model) currentMenu() menu {
	return m.menus[m.screen]
}

func (m Model) goTo(s screen) Model {
	m.screen = s
	m.cursor = 0
	m.message = ""
	return m
}
