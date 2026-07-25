package app

import (
	"DailyDashBoard/internal/tui/features"
	"DailyDashBoard/internal/tui/theme"
	"image/color"
)

// menuAmbientWords decorate menu screens when no feature is highlighted
// (the main menu, or "back" highlighted on the daily menu).
var menuAmbientWords = []string{"daily", "dashboard", "today", "curiosity", "✦"}

type menuAction int

const (
	actionOpenDailyMenu menuAction = iota
	actionQuit
	actionBack
	actionRunFeature
)

type menu struct {
	title string
	items []menuItem
}

type menuItem struct {
	label   string
	action  menuAction
	feature *features.Feature
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

// menuAmbience returns the accent color and scattered background words for
// the current menu screen: the highlighted feature's own, when the daily
// menu has one highlighted, otherwise the app's generic ambience.
func (m Model) menuAmbience() (color.Color, []string) {
	if selected, ok := m.selectedMenuItem(); ok && selected.feature != nil {
		return selected.feature.Accent, selected.feature.BackgroundWords
	}

	return theme.Primary, menuAmbientWords
}

func (m Model) goTo(s screen) Model {
	m.screen = s
	m.cursor = 0
	m.message = ""

	return m.clearContent()
}
