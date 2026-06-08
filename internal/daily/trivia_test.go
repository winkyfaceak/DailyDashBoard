package daily

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomTriviaQuestionUsesAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("amount"); got != "1" {
			t.Fatalf("expected amount query to be 1, got %q", got)
		}

		if got := r.URL.Query().Get("type"); got != "multiple" {
			t.Fatalf("expected type query to be multiple, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected JSON accept header, got %q", got)
		}

		fmt.Fprint(w, `{"response_code":0,"results":[{"category":"Science & Nature","type":"multiple","difficulty":"easy","question":"What does CPU stand for?","correct_answer":"Central Processing Unit","incorrect_answers":["Computer Personal Unit","Central Process Unit","Control Processing Unit"]}]}`)
	}))
	defer server.Close()

	restoreRandomTriviaURL(t, server.URL)

	question, err := RandomTriviaQuestion()
	if err != nil {
		t.Fatal(err)
	}

	if question.Question != "What does CPU stand for?" {
		t.Fatalf("expected question text, got %q", question.Question)
	}

	if len(question.Choices) != 4 {
		t.Fatalf("expected four choices, got %#v", question.Choices)
	}

	if !question.CheckChoice(question.CorrectChoiceIndex) {
		t.Fatal("expected correct choice index to be accepted")
	}
}

func TestDecodeRandomTriviaQuestionCleansText(t *testing.T) {
	question, err := decodeRandomTriviaQuestion(strings.NewReader(`{"response_code":0,"results":[{
		"category":"Entertainment: Music",
		"type":"multiple",
		"difficulty":"medium",
		"question":"Who wrote &quot;Purple Rain&quot;?",
		"correct_answer":"Prince",
		"incorrect_answers":["David Bowie","Michael Jackson","Stevie Wonder"]
	}]}`))
	if err != nil {
		t.Fatal(err)
	}

	if question.Question != `Who wrote "Purple Rain"?` {
		t.Fatalf("expected cleaned question, got %q", question.Question)
	}

	if question.CorrectAnswer != "Prince" {
		t.Fatalf("expected cleaned correct answer, got %q", question.CorrectAnswer)
	}
}

func TestDecodeRandomTriviaQuestionRejectsResponseCode(t *testing.T) {
	_, err := decodeRandomTriviaQuestion(strings.NewReader(`{"response_code":1,"results":[]}`))
	if err == nil {
		t.Fatal("expected response code error")
	}
}

func restoreRandomTriviaURL(t *testing.T, testURL string) {
	t.Helper()

	originalURL := randomTriviaURL
	randomTriviaURL = testURL

	t.Cleanup(func() {
		randomTriviaURL = originalURL
	})
}
