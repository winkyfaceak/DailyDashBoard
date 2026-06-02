package daily

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomTarotReadingUsesAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("n"); got != "1" {
			t.Fatalf("expected n query to be 1, got %q", got)
		}

		if got := r.URL.Query().Get("minor"); got != "true" {
			t.Fatalf("expected minor query to be true, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected JSON accept header, got %q", got)
		}

		fmt.Fprint(w, `{"nhits":1,"cards":[{"name":"The Star","name_short":"ar17","value":"17","value_int":17,"type":"major","meaning_up":"Hope and bright prospects.","meaning_rev":"Arrogance and haughtiness.","desc":"A great, radiant star."}]}`)
	}))
	defer server.Close()

	restoreRandomTarotCardURL(t, server.URL)

	reading, err := RandomTarotReading()
	if err != nil {
		t.Fatal(err)
	}

	if reading.Card.Name != "The Star" {
		t.Fatalf("expected card name, got %q", reading.Card.Name)
	}

	switch reading.Orientation {
	case TarotOrientationUpright:
		if reading.Meaning != "Hope and bright prospects." {
			t.Fatalf("expected upright meaning, got %q", reading.Meaning)
		}
	case TarotOrientationReversed:
		if reading.Meaning != "Arrogance and haughtiness." {
			t.Fatalf("expected reversed meaning, got %q", reading.Meaning)
		}
	default:
		t.Fatalf("expected tarot orientation, got %q", reading.Orientation)
	}
}

func TestDecodeRandomTarotCardCleansText(t *testing.T) {
	card, err := decodeRandomTarotCard(strings.NewReader(`{
		"cards": [{
			"name": " Ace&nbsp;of Cups ",
			"name_short": " cuac ",
			"value": " ace ",
			"value_int": 1,
			"suit": " cups ",
			"type": " minor ",
			"meaning_up": "Joy,\ncontent,\r\nabundance.",
			"meaning_rev": "Mutation&nbsp;and instability.",
			"desc": "A hand issues from a cloud."
		}]
	}`))
	if err != nil {
		t.Fatal(err)
	}

	if card.Name != "Ace of Cups" {
		t.Fatalf("expected cleaned card name, got %q", card.Name)
	}

	if card.MeaningUp != "Joy, content, abundance." {
		t.Fatalf("expected cleaned meaning, got %q", card.MeaningUp)
	}

	if card.MeaningReversed != "Mutation and instability." {
		t.Fatalf("expected cleaned reversed meaning, got %q", card.MeaningReversed)
	}
}

func TestDecodeRandomTarotCardRejectsEmptyCards(t *testing.T) {
	_, err := decodeRandomTarotCard(strings.NewReader(`{"cards":[]}`))
	if err == nil {
		t.Fatal("expected empty cards error")
	}
}

func TestTarotReadingUsesOrientationMeaning(t *testing.T) {
	card := TarotCard{
		Name:            "The Moon",
		MeaningUp:       "Hidden things rise.",
		MeaningReversed: "Confusion clears.",
	}

	reading := tarotReadingFromCard(card, TarotOrientationReversed)

	if reading.Meaning != "Confusion clears." {
		t.Fatalf("expected reversed meaning, got %q", reading.Meaning)
	}
}

func restoreRandomTarotCardURL(t *testing.T, testURL string) {
	t.Helper()

	originalURL := randomTarotCardURL
	randomTarotCardURL = testURL

	t.Cleanup(func() {
		randomTarotCardURL = originalURL
	})
}
