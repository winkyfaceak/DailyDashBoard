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
		if got := r.URL.Query().Get("n"); got != "3" {
			t.Fatalf("expected n query to be 3, got %q", got)
		}

		if got := r.URL.Query().Get("minor"); got != "true" {
			t.Fatalf("expected minor query to be true, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected JSON accept header, got %q", got)
		}

		fmt.Fprint(w, `{"nhits":3,"cards":[
			{"name":"The Star","name_short":"ar17","value":"17","value_int":17,"type":"major","meaning_up":"Hope and bright prospects.","meaning_rev":"Arrogance and haughtiness.","desc":"A great, radiant star."},
			{"name":"Ace of Cups","name_short":"cuac","value":"ace","value_int":1,"suit":"cups","type":"minor","meaning_up":"Joy, content, abundance.","meaning_rev":"Mutation and instability.","desc":"A hand issues from a cloud."},
			{"name":"Nine of Wands","name_short":"wa09","value":"nine","value_int":9,"suit":"wands","type":"minor","meaning_up":"Strength in opposition.","meaning_rev":"Obstacles and adversity.","desc":"A tired figure keeps watch."}
		]}`)
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

	if reading.SpreadName != "Past / Present / Future" {
		t.Fatalf("expected spread name, got %q", reading.SpreadName)
	}

	if len(reading.Cards) != 3 {
		t.Fatalf("expected three spread cards, got %d", len(reading.Cards))
	}

	if reading.Cards[0].Position != "Past" || reading.Cards[1].Position != "Present" || reading.Cards[2].Position != "Future" {
		t.Fatalf("expected past/present/future positions, got %#v", reading.Cards)
	}

	for _, spreadCard := range reading.Cards {
		switch spreadCard.Orientation {
		case TarotOrientationUpright:
			if spreadCard.Meaning != spreadCard.Card.MeaningUp {
				t.Fatalf("expected upright meaning, got %q", spreadCard.Meaning)
			}
		case TarotOrientationReversed:
			if spreadCard.Meaning != spreadCard.Card.MeaningReversed {
				t.Fatalf("expected reversed meaning, got %q", spreadCard.Meaning)
			}
		default:
			t.Fatalf("expected tarot orientation, got %q", spreadCard.Orientation)
		}
	}

	if reading.Cards[0].Card.ImageURL == "" {
		t.Fatal("expected card image URL")
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

	if card.ImageURL != "https://commons.wikimedia.org/wiki/Special:FilePath/Cups01.jpg?width=360" {
		t.Fatalf("expected image URL, got %q", card.ImageURL)
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

func TestTarotCardImageURLUsesCommonsFiles(t *testing.T) {
	tests := []struct {
		name string
		card TarotCard
		want string
	}{
		{
			name: "major lovers",
			card: TarotCard{Type: "major", ValueInt: 6},
			want: "https://commons.wikimedia.org/wiki/Special:FilePath/RWS_Tarot_06_Lovers.jpg?width=360",
		},
		{
			name: "minor wands nine",
			card: TarotCard{Type: "minor", Suit: "wands", ValueInt: 9},
			want: "https://commons.wikimedia.org/wiki/Special:FilePath/Tarot_Nine_of_Wands.jpg?width=360",
		},
		{
			name: "minor court fallback",
			card: TarotCard{Type: "minor", Suit: "pentacles", Value: "queen"},
			want: "https://commons.wikimedia.org/wiki/Special:FilePath/Pents13.jpg?width=360",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := TarotCardImageURL(test.card); got != test.want {
				t.Fatalf("expected image URL %q, got %q", test.want, got)
			}
		})
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
