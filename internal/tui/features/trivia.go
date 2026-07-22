package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const triviaTextWidth = 70

// RandomTrivia poses a multiple-choice trivia question.
func RandomTrivia() *Feature {
	return &Feature{
		MenuLabel: "random trivia",
		Loading:   "Fetching a trivia question...",
		ErrLabel:  "Trivia",
		BoxWidth:  84,
		Fetch: func() (Session, error) {
			question, err := daily.RandomTriviaQuestion()
			if err != nil {
				return nil, err
			}

			return &TriviaSession{
				Question:       question,
				SelectedChoice: -1,
			}, nil
		},
	}
}

// TriviaSession is the interactive multiple-choice quiz.
type TriviaSession struct {
	Question       daily.TriviaQuestion
	Submitted      bool
	SelectedChoice int
}

func (s *TriviaSession) View() string {
	return FormatTriviaQuestion(s.Question, s.Submitted, s.SelectedChoice)
}

func (s *TriviaSession) Help() string {
	return "press 1-4 to answer • esc to clear • ctrl+c to quit"
}

func (s *TriviaSession) HandleKey(msg tea.KeyPressMsg) KeyResponse {
	switch msg.String() {
	case "esc":
		return KeyResponse{Clear: true, Message: "Trivia question cleared."}
	}

	if s.Submitted {
		return KeyResponse{}
	}

	choiceIndex, ok := triviaChoiceIndex(msg.String())
	if !ok || choiceIndex >= len(s.Question.Choices) {
		return KeyResponse{}
	}

	s.Submitted = true
	s.SelectedChoice = choiceIndex

	return KeyResponse{Handled: true}
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

func FormatTriviaQuestion(question daily.TriviaQuestion, submitted bool, selectedChoice int) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(triviaTextWidth).Render("Random Trivia"))
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
