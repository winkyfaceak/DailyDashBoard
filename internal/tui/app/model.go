package app

import (
	"DailyDashBoard/internal/daily"
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
	actionRandomMath
	actionRandomTrivia
	actionDailyTarot
	actionRandomArtwork
	actionDailyNews
	actionRandomMusic
	actionMusicTrivia
	actionRandomRecipe
	actionBack
)

type Model struct {
	screen      screen
	menus       map[screen]menu
	cursor      int
	message     string
	content     contentState
	links       linkState
	audio       audioState
	math        mathState
	trivia      triviaState
	musicTrivia musicTriviaState
}

type contentState struct {
	Fact        string
	Math        string
	Trivia      string
	Dictionary  string
	Tarot       asyncContent
	Artwork     asyncContent
	News        string
	Music       asyncContent
	MusicTrivia string
	Recipe      asyncContent
}

type asyncContent struct {
	Text    string
	Pending string
}

type linkState struct {
	TarotImages  []string
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

type mathState struct {
	Active      bool
	Question    daily.MathQuestion
	AnswerInput string
	Submitted   bool
	Correct     bool
}

type triviaState struct {
	Active         bool
	Question       daily.TriviaQuestion
	Submitted      bool
	SelectedChoice int
}

type musicTriviaState struct {
	Active     bool
	Question   daily.MusicTriviaQuestion
	GuessInput string
	GuessCount int
	HintCount  int
	Solved     bool
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
					{label: "random math question", action: actionRandomMath},
					{label: "random trivia", action: actionRandomTrivia},
					{label: "daily tarot", action: actionDailyTarot},
					{label: "random artwork", action: actionRandomArtwork},
					{label: "daily news", action: actionDailyNews},
					{label: "random music", action: actionRandomMusic},
					{label: "music trivia", action: actionMusicTrivia},
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
