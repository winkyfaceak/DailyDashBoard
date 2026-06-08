package app

import (
	"DailyDashBoard/internal/tui/audio"
	"DailyDashBoard/internal/tui/features"
	"DailyDashBoard/internal/tui/platform"
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type externalOpenMsg struct {
	Err error
}

func openExternalURLCmd(targetURL string) tea.Cmd {
	return func() tea.Msg {
		return externalOpenMsg{
			Err: platform.OpenURL(targetURL),
		}
	}
}

func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "ctrl+c":
		m.stopAudio()
		return m, tea.Quit
	}

	if m.trivia.Active {
		return m.handleTriviaKeyPress(msg)
	}

	if m.math.Active {
		return m.handleMathKeyPress(msg)
	}

	if m.musicTrivia.Active {
		return m.handleMusicTriviaKeyPress(msg)
	}

	if index, ok := features.NewsLinkIndex(key); ok && len(m.links.NewsStories) > 0 {
		return m.openNewsStory(index)
	}

	switch key {
	case "q":
		m.stopAudio()
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

	case "p":
		return m.toggleAudio()

	case "o":
		return m.openPrimaryImage()

	case "r":
		return m.openRecipeSource()

	case "v":
		return m.openRecipeVideo()

	case "s":
		return m.openMusicTrack()

	case "enter":
		selected, ok := m.selectedMenuItem()
		if !ok {
			return m, nil
		}

		return m.handleSelection(selected)
	}

	return m, nil
}

func (m Model) openNewsStory(index int) (tea.Model, tea.Cmd) {
	if index >= len(m.links.NewsStories) || m.links.NewsStories[index] == "" {
		m.message = fmt.Sprintf("Story %d does not include a link.", index+1)
		return m, nil
	}

	m.message = fmt.Sprintf("Opening story %d...", index+1)
	return m, openExternalURLCmd(m.links.NewsStories[index])
}

func (m Model) toggleAudio() (tea.Model, tea.Cmd) {
	if m.audio.URL == "" {
		return m, nil
	}

	audioLabel := m.currentAudioLabel()

	if m.audio.Player != nil {
		if m.audio.Paused {
			if err := m.audio.Player.Resume(); err != nil {
				m.message = fmt.Sprintf("Audio error: %v", err)
				return m, nil
			}

			m.audio.Paused = false
			m.message = "Resumed " + audioLabel + "."
			return m, nil
		}

		if err := m.audio.Player.Pause(); err != nil {
			m.message = fmt.Sprintf("Audio error: %v", err)
			return m, nil
		}

		m.audio.Paused = true
		m.message = "Paused " + audioLabel + "."
		return m, nil
	}

	if m.audio.Starting {
		m.message = "Starting " + audioLabel + "..."
		return m, nil
	}

	m.message = "Playing " + audioLabel + "..."
	m.audio.Starting = true
	m.audio.RequestID++

	return m, audio.StartCmd(m.audio.URL, m.audio.AllowDownload, m.audio.RequestID)
}

func (m Model) openPrimaryImage() (tea.Model, tea.Cmd) {
	if len(m.links.TarotImages) > 0 {
		m.message = "Opening first tarot card image..."
		return m, openExternalURLCmd(m.links.TarotImages[0])
	}

	if m.links.ArtworkImage != "" {
		m.message = "Opening full artwork image..."
		return m, openExternalURLCmd(m.links.ArtworkImage)
	}

	if m.links.MusicImage != "" {
		m.message = "Opening album art..."
		return m, openExternalURLCmd(m.links.MusicImage)
	}

	if m.links.RecipeImage != "" {
		m.message = "Opening food image..."
		return m, openExternalURLCmd(m.links.RecipeImage)
	}

	return m, nil
}

func (m Model) openRecipeSource() (tea.Model, tea.Cmd) {
	if m.links.RecipeSource == "" {
		return m, nil
	}

	m.message = "Opening recipe source..."
	return m, openExternalURLCmd(m.links.RecipeSource)
}

func (m Model) openRecipeVideo() (tea.Model, tea.Cmd) {
	if m.links.RecipeVideo == "" {
		return m, nil
	}

	m.message = "Opening cooking video..."
	return m, openExternalURLCmd(m.links.RecipeVideo)
}

func (m Model) openMusicTrack() (tea.Model, tea.Cmd) {
	if m.links.MusicTrack == "" {
		return m, nil
	}

	m.message = "Opening song..."
	return m, openExternalURLCmd(m.links.MusicTrack)
}

func (m Model) handleTriviaKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m = m.clearResultContent()
		m.message = "Trivia question cleared."
		return m, nil
	}

	if m.trivia.Submitted {
		return m, nil
	}

	choiceIndex, ok := triviaChoiceIndex(msg.String())
	if !ok || choiceIndex >= len(m.trivia.Question.Choices) {
		return m, nil
	}

	m.message = ""
	m.trivia.Submitted = true
	m.trivia.SelectedChoice = choiceIndex
	m.content.Trivia = m.formatTriviaContent()

	return m, nil
}

func triviaChoiceIndex(key string) (int, bool) {
	switch key {
	case "1":
		return 0, true
	case "2":
		return 1, true
	case "3":
		return 2, true
	case "4":
		return 3, true
	default:
		return 0, false
	}
}

func (m Model) handleMathKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m.submitMathAnswer()
	case "backspace", "ctrl+h":
		if m.math.AnswerInput != "" {
			m.math.AnswerInput = m.math.AnswerInput[:len(m.math.AnswerInput)-1]
			m.math.Submitted = false
			m.math.Correct = false
			m.message = ""
			m.content.Math = m.formatMathContent()
		}

		return m, nil
	case "esc":
		m = m.clearResultContent()
		m.message = "Math question cleared."
		return m, nil
	}

	return m.appendMathAnswerText(msg.Key().Text), nil
}

func (m Model) handleMusicTriviaKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m.submitMusicTriviaGuess()
	case "backspace", "ctrl+h":
		if m.musicTrivia.GuessInput != "" {
			runes := []rune(m.musicTrivia.GuessInput)
			m.musicTrivia.GuessInput = string(runes[:len(runes)-1])
			m.message = ""
			m.content.MusicTrivia = m.formatMusicTriviaContent()
		}

		return m, nil
	case "esc":
		m.stopAudio()
		m = m.clearResultContent()
		m.message = "Music trivia cleared."
		return m, nil
	case "ctrl+p":
		return m.toggleAudio()
	case "?":
		return m.revealMusicTriviaHint()
	}

	return m.appendMusicTriviaGuessText(msg.Key().Text), nil
}

func (m Model) appendMusicTriviaGuessText(text string) Model {
	if text == "" || m.musicTrivia.Solved {
		return m
	}

	for _, char := range text {
		if char == '\n' || char == '\r' || char == '\t' {
			continue
		}

		if len([]rune(m.musicTrivia.GuessInput)) >= 80 {
			break
		}

		m.musicTrivia.GuessInput += string(char)
	}

	m.message = ""
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	return m
}

func (m Model) submitMusicTriviaGuess() (tea.Model, tea.Cmd) {
	guess := strings.TrimSpace(m.musicTrivia.GuessInput)
	if guess == "" {
		m.message = "Enter a guess first."
		return m, nil
	}

	if m.musicTrivia.Question.CheckGuess(guess) {
		m.message = ""
		m.musicTrivia.Solved = true
		m.content.MusicTrivia = m.formatMusicTriviaContent()
		return m, nil
	}

	m.musicTrivia.GuessCount++
	m.musicTrivia.GuessInput = ""
	m.message = "Not quite. Try again."
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	return m, nil
}

func (m Model) revealMusicTriviaHint() (tea.Model, tea.Cmd) {
	hints := m.musicTrivia.Question.Hints()
	if m.musicTrivia.HintCount >= len(hints) {
		m.message = "No more hints available."
		return m, nil
	}

	neededGuesses := (m.musicTrivia.HintCount + 1) * 2
	if m.musicTrivia.GuessCount < neededGuesses {
		m.message = fmt.Sprintf("Hints unlock after %d wrong guesses.", neededGuesses)
		return m, nil
	}

	m.message = ""
	m.musicTrivia.HintCount++
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	return m, nil
}

func (m Model) appendMathAnswerText(text string) Model {
	if text == "" {
		return m
	}

	if m.math.Submitted {
		m.math.AnswerInput = ""
		m.math.Submitted = false
		m.math.Correct = false
	}

	for _, char := range text {
		if char < '0' || char > '9' {
			continue
		}

		if len(m.math.AnswerInput) >= 8 {
			break
		}

		m.math.AnswerInput += string(char)
	}

	m.message = ""
	m.content.Math = m.formatMathContent()

	return m
}

func (m Model) submitMathAnswer() (tea.Model, tea.Cmd) {
	answerInput := strings.TrimSpace(m.math.AnswerInput)
	if answerInput == "" {
		m.message = "Enter an answer first."
		return m, nil
	}

	answer, err := strconv.Atoi(answerInput)
	if err != nil {
		m.message = "Use whole numbers only."
		return m, nil
	}

	m.message = ""
	m.math.Submitted = true
	m.math.Correct = m.math.Question.CheckAnswer(answer)
	m.content.Math = m.formatMathContent()

	return m, nil
}

func (m Model) formatMathContent() string {
	return features.FormatMathQuestion(
		m.math.Question,
		m.math.AnswerInput,
		m.math.Submitted,
		m.math.Correct,
	)
}

func (m Model) formatTriviaContent() string {
	return features.FormatTriviaQuestion(
		m.trivia.Question,
		m.trivia.Submitted,
		m.trivia.SelectedChoice,
	)
}

func (m Model) formatMusicTriviaContent() string {
	return features.FormatMusicTriviaQuestion(
		m.musicTrivia.Question,
		m.musicTrivia.GuessInput,
		m.musicTrivia.GuessCount,
		m.musicTrivia.HintCount,
		m.musicTrivia.Solved,
	)
}

func (m Model) handleSelection(selected menuItem) (tea.Model, tea.Cmd) {
	switch selected.action {
	case actionOpenDailyMenu:
		return m.goTo(dailyMenu), nil
	case actionQuit:
		m.stopAudio()
		return m, tea.Quit
	case actionRandomWord:
		return m.startFeature("Searching for a random word...", features.FetchRandomDictionaryCmd("en"))
	case actionRandomFact:
		return m.startFeature("Fetching a random fact...", features.FetchFactCmd())
	case actionRandomMath:
		return m.startFeature("Writing a random math question...", features.FetchMathCmd())
	case actionRandomTrivia:
		return m.startFeature("Fetching a trivia question...", features.FetchTriviaCmd())
	case actionDailyTarot:
		return m.startFeature("Drawing a three-card tarot spread...", features.FetchTarotCmd())
	case actionRandomArtwork:
		return m.startFeature("Fetching artwork...", features.FetchArtworkCmd())
	case actionDailyNews:
		return m.startFeature("Fetching daily news...", features.FetchNewsCmd())
	case actionRandomMusic:
		return m.startFeature("Finding a genre and song...", features.FetchMusicCmd())
	case actionMusicTrivia:
		return m.startFeature("Finding a mystery music preview...", features.FetchMusicTriviaCmd())
	case actionRandomRecipe:
		return m.startFeature("Finding a random recipe...", features.FetchRecipeCmd())
	case actionBack:
		return m.goTo(mainMenu), nil
	default:
		m = m.clearContent()
		m.message = fmt.Sprintf("No action is configured for %q.", selected.label)
		return m, nil
	}
}

func (m Model) startFeature(message string, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m = m.clearContent()
	m.message = message

	return m, cmd
}
