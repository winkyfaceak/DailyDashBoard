package features

import (
	"DailyDashBoard/internal/daily"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const musicTriviaTextWidth = 66

type MusicTriviaResultMsg struct {
	Question daily.MusicTriviaQuestion
	Err      error
}

func FetchMusicTriviaCmd() tea.Cmd {
	return func() tea.Msg {
		question, err := daily.RandomMusicTriviaQuestion()

		return MusicTriviaResultMsg{
			Question: question,
			Err:      err,
		}
	}
}

func FormatMusicTriviaQuestion(question daily.MusicTriviaQuestion, guessInput string, guessCount int, hintCount int, solved bool) string {
	var builder strings.Builder

	builder.WriteString(musicTriviaHeaderStyle.Render("Music Trivia"))
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
