package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const triviaTextWidth = 70

type TriviaResultMsg struct {
	Question daily.TriviaQuestion
	Err      error
}

func FetchTriviaCmd() tea.Cmd {
	return func() tea.Msg {
		question, err := daily.RandomTriviaQuestion()

		return TriviaResultMsg{
			Question: question,
			Err:      err,
		}
	}
}

func FormatTriviaQuestion(question daily.TriviaQuestion, submitted bool, selectedChoice int) string {
	var builder strings.Builder

	builder.WriteString(triviaHeaderStyle.Render("Random Trivia"))
	builder.WriteString("\n\n")

	meta := formatTriviaMeta(question)
	if meta != "" {
		builder.WriteString(subtleTextStyle.Render(meta))
		builder.WriteString("\n\n")
	}

	builder.WriteString(textutil.Wrap(question.Question, triviaTextWidth))
	builder.WriteString("\n\n")

	for index, choice := range question.Choices {
		line := fmt.Sprintf("%d. %s", index+1, choice)
		if submitted {
			switch {
			case index == question.CorrectChoiceIndex:
				line = successTextStyle.Render(line)
			case index == selectedChoice:
				line = subtleTextStyle.Render(line)
			}
		}

		builder.WriteString(line)
		if index < len(question.Choices)-1 {
			builder.WriteString("\n")
		}
	}

	if submitted {
		builder.WriteString("\n\n")
		if question.CheckChoice(selectedChoice) {
			builder.WriteString(successTextStyle.Render("Correct."))
		} else {
			builder.WriteString(subtleTextStyle.Render(fmt.Sprintf("Answer: %s", question.CorrectAnswer)))
		}
	}

	return builder.String()
}

func formatTriviaMeta(question daily.TriviaQuestion) string {
	var parts []string
	if question.Category != "" {
		parts = append(parts, question.Category)
	}

	if question.Difficulty != "" {
		parts = append(parts, titleCase(question.Difficulty))
	}

	return strings.Join(parts, " • ")
}
