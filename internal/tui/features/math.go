package features

import (
	"DailyDashBoard/internal/daily"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const mathTextWidth = 66

type MathResultMsg struct {
	Question daily.MathQuestion
	Err      error
}

func FetchMathCmd() tea.Cmd {
	return func() tea.Msg {
		question, err := daily.RandomMathQuestion()

		return MathResultMsg{
			Question: question,
			Err:      err,
		}
	}
}

func FormatMathQuestion(question daily.MathQuestion, answerInput string, submitted bool, correct bool) string {
	var builder strings.Builder

	builder.WriteString(mathHeaderStyle.Render("Random Math"))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("Question"))
	builder.WriteString("\n")
	builder.WriteString(mathQuestionStyle.Render(question.Prompt()))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("Answer"))
	builder.WriteString("\n")
	if answerInput == "" {
		builder.WriteString(exampleStyle.Render("> _"))
	} else if submitted {
		builder.WriteString(fmt.Sprintf("> %s", answerInput))
	} else {
		builder.WriteString(fmt.Sprintf("> %s_", answerInput))
	}

	if submitted {
		builder.WriteString("\n\n")
		if correct {
			builder.WriteString(successTextStyle.Render("Correct. Nice work."))
		} else {
			builder.WriteString(subtleTextStyle.Render(fmt.Sprintf("Not quite. The answer is %d.", question.Answer)))
		}
	}

	return builder.String()
}
