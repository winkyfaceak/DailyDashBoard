package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const (
	mathTextWidth      = 66
	mathAnswerMaxRunes = 8
)

// RandomMath poses a random arithmetic question.
func RandomMath() *Feature {
	return &Feature{
		MenuLabel: "random math question",
		Loading:   "Writing a random math question...",
		ErrLabel:  "Math",
		BoxWidth:  78,
		Fetch: func() (Session, error) {
			question, err := daily.RandomMathQuestion()
			if err != nil {
				return nil, err
			}

			return &MathSession{Question: question}, nil
		},
	}
}

// MathSession is the interactive math quiz.
type MathSession struct {
	Question    daily.MathQuestion
	AnswerInput string
	Submitted   bool
	Correct     bool
}

func (s *MathSession) View() string {
	return FormatMathQuestion(s.Question, s.AnswerInput, s.Submitted, s.Correct)
}

func (s *MathSession) Help() string {
	return "type answer • enter to check • backspace to edit • esc to clear • ctrl+c to quit"
}

func (s *MathSession) HandleKey(msg tea.KeyPressMsg) KeyResponse {
	switch msg.String() {
	case "enter":
		return s.submit()
	case "backspace", "ctrl+h":
		if s.AnswerInput == "" {
			return KeyResponse{}
		}

		s.AnswerInput = trimLastRune(s.AnswerInput)
		s.Submitted = false
		s.Correct = false

		return KeyResponse{Handled: true}
	case "esc":
		return KeyResponse{Clear: true, Message: "Math question cleared."}
	}

	return s.appendAnswerText(msg.Key().Text)
}

func (s *MathSession) appendAnswerText(text string) KeyResponse {
	if text == "" {
		return KeyResponse{}
	}

	if s.Submitted {
		s.AnswerInput = ""
		s.Submitted = false
		s.Correct = false
	}

	s.AnswerInput = appendFiltered(s.AnswerInput, text, mathAnswerMaxRunes, func(char rune) bool {
		return char >= '0' && char <= '9'
	})

	return KeyResponse{Handled: true}
}

func (s *MathSession) submit() KeyResponse {
	answerInput := strings.TrimSpace(s.AnswerInput)
	if answerInput == "" {
		return KeyResponse{Handled: true, Message: "Enter an answer first."}
	}

	answer, err := strconv.Atoi(answerInput)
	if err != nil {
		return KeyResponse{Handled: true, Message: "Use whole numbers only."}
	}

	s.Submitted = true
	s.Correct = s.Question.CheckAnswer(answer)

	return KeyResponse{Handled: true}
}

func FormatMathQuestion(question daily.MathQuestion, answerInput string, submitted bool, correct bool) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(mathTextWidth).Render("Random Math"))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("Question"))
	builder.WriteString("\n")
	builder.WriteString(strongTextStyle.Render(question.Prompt()))
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
