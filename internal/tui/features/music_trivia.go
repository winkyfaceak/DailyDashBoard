package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	musicTriviaTextWidth     = 66
	musicTriviaGuessMaxRunes = 80
)

// MusicTrivia plays a mystery song preview and asks who or what it is.
func MusicTrivia() *Feature {
	return &Feature{
		MenuLabel:       "music trivia",
		Loading:         "Finding a mystery music preview...",
		ErrLabel:        "Music trivia",
		BoxWidth:        84,
		Accent:          lipgloss.Color("#8C7AE6"),
		BackgroundWords: []string{"name that tune", "beat", "lyric", "♪", "?"},
		Fetch: func() (Session, error) {
			question, err := daily.RandomMusicTriviaQuestion()
			if err != nil {
				return nil, err
			}

			return &MusicTriviaSession{Question: question}, nil
		},
	}
}

// MusicTriviaSession is the interactive guess-the-song quiz.
type MusicTriviaSession struct {
	Question   daily.MusicTriviaQuestion
	GuessInput string
	GuessCount int
	HintCount  int
	Solved     bool
}

func (s *MusicTriviaSession) View() string {
	return FormatMusicTriviaQuestion(s.Question, s.GuessInput, s.GuessCount, s.HintCount, s.Solved)
}

func (s *MusicTriviaSession) Audio() Audio {
	return Audio{
		URL:           s.Question.Pick.PreviewURL,
		Label:         "mystery preview",
		AllowDownload: true,
	}
}

func (s *MusicTriviaSession) Help() string {
	return "type guess • enter to submit • ctrl+p to play preview • ? for hint • backspace to edit • esc to menu • ctrl+c to quit"
}

func (s *MusicTriviaSession) HandleKey(msg tea.KeyPressMsg) KeyResponse {
	switch msg.String() {
	case "enter":
		return s.submitGuess()
	case "backspace", "ctrl+h":
		if s.GuessInput == "" {
			return KeyResponse{}
		}

		s.GuessInput = trimLastRune(s.GuessInput)

		return KeyResponse{Handled: true}
	case "esc":
		return KeyResponse{Clear: true, Message: "Music trivia cleared."}
	case "ctrl+p":
		return KeyResponse{ToggleAudio: true}
	case "?":
		return s.revealHint()
	}

	return s.appendGuessText(msg.Key().Text)
}

func (s *MusicTriviaSession) appendGuessText(text string) KeyResponse {
	if text == "" || s.Solved {
		return KeyResponse{}
	}

	s.GuessInput = appendFiltered(s.GuessInput, text, musicTriviaGuessMaxRunes, func(char rune) bool {
		return char != '\n' && char != '\r' && char != '\t'
	})

	return KeyResponse{Handled: true}
}

func (s *MusicTriviaSession) submitGuess() KeyResponse {
	guess := strings.TrimSpace(s.GuessInput)
	if guess == "" {
		return KeyResponse{Handled: true, Message: "Enter a guess first."}
	}

	if s.Question.CheckGuess(guess) {
		s.Solved = true
		return KeyResponse{Handled: true}
	}

	s.GuessCount++
	s.GuessInput = ""

	return KeyResponse{Handled: true, Message: "Not quite. Try again."}
}

func (s *MusicTriviaSession) revealHint() KeyResponse {
	hints := s.Question.Hints()
	if s.HintCount >= len(hints) {
		return KeyResponse{Handled: true, Message: "No more hints available."}
	}

	neededGuesses := (s.HintCount + 1) * 2
	if s.GuessCount < neededGuesses {
		return KeyResponse{Handled: true, Message: fmt.Sprintf("Hints unlock after %d wrong guesses.", neededGuesses)}
	}

	s.HintCount++

	return KeyResponse{Handled: true}
}

func FormatMusicTriviaQuestion(question daily.MusicTriviaQuestion, guessInput string, guessCount int, hintCount int, solved bool) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(musicTriviaTextWidth).Render("Music Trivia"))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("Mystery"))
	builder.WriteString("\n")
	builder.WriteString(question.Prompt())
	builder.WriteString("\n")
	builder.WriteString(exampleStyle.Render("Press ctrl+p to play the preview."))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("Guess"))
	builder.WriteString("\n")
	if guessInput == "" {
		builder.WriteString(exampleStyle.Render("> _"))
	} else if solved {
		builder.WriteString(fmt.Sprintf("> %s", guessInput))
	} else {
		builder.WriteString(fmt.Sprintf("> %s_", guessInput))
	}

	if guessCount > 0 && !solved {
		builder.WriteString("\n")
		builder.WriteString(subtleTextStyle.Render(fmt.Sprintf("Guesses: %d", guessCount)))
	}

	if hintCount > 0 {
		builder.WriteString("\n\n")
		builder.WriteString(labelStyle.Render("Hints"))
		builder.WriteString("\n")

		hints := question.Hints()
		if hintCount > len(hints) {
			hintCount = len(hints)
		}

		for index, hint := range hints[:hintCount] {
			builder.WriteString(subtleTextStyle.Render(fmt.Sprintf("%d. %s", index+1, hint)))
			if index < hintCount-1 {
				builder.WriteString("\n")
			}
		}
	}

	if solved {
		builder.WriteString("\n\n")
		builder.WriteString(successTextStyle.Render(fmt.Sprintf("Correct: %s", formatMusicTriviaAnswer(question))))
	} else {
		builder.WriteString("\n\n")
		builder.WriteString(subtleTextStyle.Render(formatMusicTriviaHintStatus(question, guessCount, hintCount)))
	}

	return builder.String()
}

func formatMusicTriviaHintStatus(question daily.MusicTriviaQuestion, guessCount int, hintCount int) string {
	if hintCount >= len(question.Hints()) {
		return "No more hints available."
	}

	neededGuesses := (hintCount + 1) * 2
	if guessCount >= neededGuesses {
		return "Press ? for a hint."
	}

	return fmt.Sprintf("Hints unlock after %d wrong guesses.", neededGuesses)
}

func formatMusicTriviaAnswer(question daily.MusicTriviaQuestion) string {
	parts := []string{question.Pick.TrackName, "by", question.Pick.ArtistName}
	answer := strings.Join(parts, " ")
	if question.Pick.AlbumName != "" {
		answer += " (" + question.Pick.AlbumName + ")"
	}

	return answer
}
