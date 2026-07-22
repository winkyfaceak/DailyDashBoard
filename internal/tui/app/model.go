// Package app is the generic bubbletea shell of the dashboard. It renders
// menus, routes messages, and hosts whichever feature session is active; it
// never references individual features directly.
package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"

	tea "charm.land/bubbletea/v2"
)

type screen int

const (
	mainMenu screen = iota
	dailyMenu
)

type Model struct {
	screen  screen
	menus   map[screen]menu
	cursor  int
	message string

	// feature and session describe the currently displayed result, if any.
	feature *features.Feature
	session features.Session

	audio audioState
}

type audioState struct {
	URL           string
	Label         string
	AllowDownload bool
	Player        *audio.Player
	Paused        bool
	Starting      bool
	RequestID     int
}

func InitialModel() Model {
	registry := features.Daily()

	dailyItems := make([]menuItem, 0, len(registry)+1)
	for _, feature := range registry {
		dailyItems = append(dailyItems, menuItem{
			label:   feature.MenuLabel,
			action:  actionRunFeature,
			feature: feature,
		})
	}

	dailyItems = append(dailyItems, menuItem{label: "back", action: actionBack})

	return Model{
		screen: mainMenu,
		menus: map[screen]menu{
			mainMenu: {
				title: "Daily Dashboard",
				items: []menuItem{
					{label: "dailies", action: actionOpenDailyMenu},
					{label: "exit", action: actionQuit},
				},
			},
			dailyMenu: {
				title: "Daily Picks",
				items: dailyItems,
			},
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
