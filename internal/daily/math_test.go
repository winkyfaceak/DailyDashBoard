package daily

import "testing"

func TestRandomMathQuestionHasIntegerAnswer(t *testing.T) {
	for range 100 {
		question, err := RandomMathQuestion()
		if err != nil {
			t.Fatal(err)
		}

		if question.Operator == "" {
			t.Fatal("expected operator")
		}

		if question.Prompt() == "" {
			t.Fatal("expected prompt")
		}

		if !question.CheckAnswer(question.Answer) {
			t.Fatalf("expected answer %d to be accepted for %#v", question.Answer, question)
		}

		switch question.Operator {
		case "+":
			if question.Left+question.Right != question.Answer {
				t.Fatalf("expected addition answer, got %#v", question)
			}
		case "-":
			if question.Left-question.Right != question.Answer {
				t.Fatalf("expected subtraction answer, got %#v", question)
			}

			if question.Answer < 0 {
				t.Fatalf("expected non-negative subtraction answer, got %#v", question)
			}
		case "x":
			if question.Left*question.Right != question.Answer {
				t.Fatalf("expected multiplication answer, got %#v", question)
			}
		case "/":
			if question.Right == 0 {
				t.Fatal("division question used zero divisor")
			}

			if question.Left/question.Right != question.Answer || question.Left%question.Right != 0 {
				t.Fatalf("expected exact division answer, got %#v", question)
			}
		default:
			t.Fatalf("unexpected operator %q", question.Operator)
		}
	}
}

func TestMathQuestionCheckAnswerRejectsWrongAnswer(t *testing.T) {
	question := MathQuestion{Answer: 42}

	if question.CheckAnswer(41) {
		t.Fatal("expected wrong answer to be rejected")
	}
}

func TestRandomMathNumberRejectsInvalidRange(t *testing.T) {
	_, err := randomMathNumber(10, 1)
	if err == nil {
		t.Fatal("expected invalid range error")
	}
}
