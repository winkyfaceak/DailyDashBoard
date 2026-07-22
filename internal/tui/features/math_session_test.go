package features

import (
	"DailyDashBoard/internal/daily"
	"strings"
	"testing"
)

func TestMathAnswerFlow(t *testing.T) {
	s := &MathSession{
		Question: daily.MathQuestion{Left: 6, Right: 7, Operator: "x", Answer: 42},
	}

	s.HandleKey(keyPress("4"))
	s.HandleKey(keyPress("2"))
	s.HandleKey(keyPress("enter"))

	if !s.Submitted {
		t.Fatal("expected answer to be submitted")
	}

	if !s.Correct {
		t.Fatal("expected answer to be correct")
	}

	if !strings.Contains(s.View(), "Correct") {
		t.Fatalf("expected correct feedback, got %q", s.View())
	}
}

func TestMathBackspaceEditsAnswer(t *testing.T) {
	s := &MathSession{
		Question:    daily.MathQuestion{Left: 9, Right: 3, Operator: "+", Answer: 12},
		AnswerInput: "123",
		Submitted:   true,
	}

	s.HandleKey(keyPress("backspace"))

	if s.AnswerInput != "12" {
		t.Fatalf("expected edited answer input, got %q", s.AnswerInput)
	}

	if s.Submitted {
		t.Fatal("expected edit to clear submitted state")
	}
}

func TestMathEscapeClearsQuestion(t *testing.T) {
	s := &MathSession{
		Question: daily.MathQuestion{Left: 9, Right: 3, Operator: "+", Answer: 12},
	}

	response := s.HandleKey(keyPress("esc"))

	if !response.Clear {
		t.Fatal("expected escape to clear the session")
	}

	if response.Message != "Math question cleared." {
		t.Fatalf("expected clear message, got %q", response.Message)
	}
}

func TestMathIgnoresNonDigits(t *testing.T) {
	s := &MathSession{
		Question: daily.MathQuestion{Left: 9, Right: 3, Operator: "+", Answer: 12},
	}

	s.HandleKey(keyPress("a"))
	s.HandleKey(keyPress("1"))

	if s.AnswerInput != "1" {
		t.Fatalf("expected only digits to be kept, got %q", s.AnswerInput)
	}
}
