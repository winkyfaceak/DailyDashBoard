package app

import (
	"DailyDashBoard/internal/tui/audio"

	tea "charm.land/bubbletea/v2"
)

type screen int

const (
	mainMenu screen = iota
	dailyMenu
)

type menuAction int

const (
	actionOpenDailyMenu menuAction = iota
	actionQuit
	actionRandomWord
	actionRandomFact
	actionRandomArtwork
	actionDailyNews
	actionRandomMusic
	actionRandomRecipe
	actionBack
)

type Model struct {
	screen  screen
	menus   map[screen]menu
	cursor  int
	message string
	content contentState
	links   linkState
	audio   audioState
}

type contentState struct {
	Fact       string
	Dictionary string
	Artwork    asyncContent
	News       string
	Music      asyncContent
	Recipe     asyncContent
}

type asyncContent struct {
	Text    string
	Pending string
}

type linkState struct {
	ArtworkImage string
	MusicImage   string
	MusicTrack   string
	NewsStories  []string
	RecipeImage  string
	RecipeSource string
	RecipeVideo  string
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
				items: []menuItem{
					{label: "random word", action: actionRandomWord},
					{label: "random fact", action: actionRandomFact},
					{label: "random artwork", action: actionRandomArtwork},
					{label: "daily news", action: actionDailyNews},
					{label: "random music", action: actionRandomMusic},
					{label: "random recipe", action: actionRandomRecipe},
					{label: "back", action: actionBack},
				},
			},
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
