package features

import (
	"DailyDashBoard/internal/daily"
	"strings"
	"testing"
)

func testTriviaSession() *TriviaSession {
	return &TriviaSession{
		Question: daily.TriviaQuestion{
			Question:           "Which planet is known as the Red Planet?",
			CorrectAnswer:      "Mars",
			Choices:            []string{"Venus", "Mars", "Jupiter", "Saturn"},
			CorrectChoiceIndex: 1,
		},
		SelectedChoice: -1,
	}
}

func TestTriviaChoiceFlow(t *testing.T) {
	s := testTriviaSession()

	s.HandleKey(keyPress("2"))

	if !s.Submitted {
		t.Fatal("expected trivia answer to be submitted")
	}

	if s.SelectedChoice != 1 {
		t.Fatalf("expected selected choice 1, got %d", s.SelectedChoice)
	}

	if !strings.Contains(s.View(), "Correct") {
		t.Fatalf("expected correct feedback, got %q", s.View())
	}
}

func TestTriviaEscapeClearsQuestion(t *testing.T) {
	s := testTriviaSession()

	response := s.HandleKey(keyPress("esc"))

	if !response.Clear {
		t.Fatal("expected escape to clear the session")
	}

	if response.Message != "Trivia question cleared." {
		t.Fatalf("expected clear message, got %q", response.Message)
	}
}

func TestTriviaIgnoresChoicesAfterSubmit(t *testing.T) {
	s := testTriviaSession()

	s.HandleKey(keyPress("2"))
	response := s.HandleKey(keyPress("3"))

	if response.Handled {
		t.Fatal("expected key after submit to be ignored")
	}

	if s.SelectedChoice != 1 {
		t.Fatalf("expected selected choice to stay 1, got %d", s.SelectedChoice)
	}
}
