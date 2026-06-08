package app

import (
	"DailyDashBoard/internal/daily"
	"strings"
	"testing"
)

func TestTriviaChoiceFlow(t *testing.T) {
	m := InitialModel()
	m.trivia = triviaState{
		Active: true,
		Question: daily.TriviaQuestion{
			Question:           "Which planet is known as the Red Planet?",
			CorrectAnswer:      "Mars",
			Choices:            []string{"Venus", "Mars", "Jupiter", "Saturn"},
			CorrectChoiceIndex: 1,
		},
		SelectedChoice: -1,
	}
	m.content.Trivia = m.formatTriviaContent()

	updated, _ := m.handleTriviaKeyPress(keyPress("2"))
	m = updated.(Model)

	if !m.trivia.Submitted {
		t.Fatal("expected trivia answer to be submitted")
	}

	if m.trivia.SelectedChoice != 1 {
		t.Fatalf("expected selected choice 1, got %d", m.trivia.SelectedChoice)
	}

	if !strings.Contains(m.content.Trivia, "Correct") {
		t.Fatalf("expected correct feedback, got %q", m.content.Trivia)
	}
}

func TestTriviaEscapeClearsQuestion(t *testing.T) {
	m := InitialModel()
	m.trivia = triviaState{
		Active: true,
		Question: daily.TriviaQuestion{
			Question:           "Which planet is known as the Red Planet?",
			CorrectAnswer:      "Mars",
			Choices:            []string{"Venus", "Mars", "Jupiter", "Saturn"},
			CorrectChoiceIndex: 1,
		},
	}
	m.content.Trivia = m.formatTriviaContent()

	updated, _ := m.handleTriviaKeyPress(keyPress("esc"))
	m = updated.(Model)

	if m.trivia.Active {
		t.Fatal("expected trivia question to be inactive")
	}

	if m.content.Trivia != "" {
		t.Fatalf("expected trivia content to be cleared, got %q", m.content.Trivia)
	}
}
