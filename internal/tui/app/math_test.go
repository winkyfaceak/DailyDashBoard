package app

import (
	"DailyDashBoard/internal/daily"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestMathAnswerFlow(t *testing.T) {
	m := InitialModel()
	m.math = mathState{
		Active:   true,
		Question: daily.MathQuestion{Left: 6, Right: 7, Operator: "x", Answer: 42},
	}
	m.content.Math = m.formatMathContent()

	updated, _ := m.handleMathKeyPress(keyPress("4"))
	m = updated.(Model)
	updated, _ = m.handleMathKeyPress(keyPress("2"))
	m = updated.(Model)
	updated, _ = m.handleMathKeyPress(keyPress("enter"))
	m = updated.(Model)

	if !m.math.Submitted {
		t.Fatal("expected answer to be submitted")
	}

	if !m.math.Correct {
		t.Fatal("expected answer to be correct")
	}

	if !strings.Contains(m.content.Math, "Correct") {
		t.Fatalf("expected correct feedback, got %q", m.content.Math)
	}
}

func TestMathBackspaceEditsAnswer(t *testing.T) {
	m := InitialModel()
	m.math = mathState{
		Active:      true,
		Question:    daily.MathQuestion{Left: 9, Right: 3, Operator: "+", Answer: 12},
		AnswerInput: "123",
		Submitted:   true,
	}
	m.content.Math = m.formatMathContent()

	updated, _ := m.handleMathKeyPress(keyPress("backspace"))
	m = updated.(Model)

	if m.math.AnswerInput != "12" {
		t.Fatalf("expected edited answer input, got %q", m.math.AnswerInput)
	}

	if m.math.Submitted {
		t.Fatal("expected edit to clear submitted state")
	}
}

func TestMathEscapeClearsQuestion(t *testing.T) {
	m := InitialModel()
	m.math = mathState{
		Active:   true,
		Question: daily.MathQuestion{Left: 9, Right: 3, Operator: "+", Answer: 12},
	}
	m.content.Math = m.formatMathContent()

	updated, _ := m.handleMathKeyPress(keyPress("esc"))
	m = updated.(Model)

	if m.math.Active {
		t.Fatal("expected math question to be inactive")
	}

	if m.content.Math != "" {
		t.Fatalf("expected math content to be cleared, got %q", m.content.Math)
	}
}

func keyPress(value string) tea.KeyPressMsg {
	if value == "enter" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	}

	if value == "backspace" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyBackspace})
	}

	if value == "esc" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape})
	}

	return tea.KeyPressMsg(tea.Key{Text: value, Code: []rune(value)[0]})
}
