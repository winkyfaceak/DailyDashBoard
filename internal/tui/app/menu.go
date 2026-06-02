package app

type menu struct {
	title string
	items []menuItem
}

type menuItem struct {
	label  string
	action menuAction
}

func (m Model) currentMenu() menu {
	return m.menus[m.screen]
}

func (m Model) selectedMenuItem() (menuItem, bool) {
	items := m.currentMenu().items
	if m.cursor < 0 || m.cursor >= len(items) {
		return menuItem{}, false
	}

	return items[m.cursor], true
}

func (m Model) goTo(s screen) Model {
	m.screen = s
	m.cursor = 0
	m.message = ""

	return m.clearContent()
}
