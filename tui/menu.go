package tui

type menu struct {
	title string
	items []string
}

func (m model) currentMenu() menu {
	return m.menus[m.screen]
}

func (m model) goTo(s screen) model {
	m.screen = s
	m.cursor = 0
	m.message = ""
	return m
}
