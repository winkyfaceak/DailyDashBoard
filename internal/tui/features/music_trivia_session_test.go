package features

import (
	"DailyDashBoard/internal/daily"
	"testing"
)

func TestMusicTriviaGuessFlow(t *testing.T) {
	s := &MusicTriviaSession{
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName: "Prince",
				TrackName:  "Kiss",
				PreviewURL: "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessArtist,
		},
	}

	for _, key := range []string{"P", "r", "i", "n", "c", "e"} {
		s.HandleKey(keyPress(key))
	}

	s.HandleKey(keyPress("enter"))

	if !s.Solved {
		t.Fatal("expected music trivia to be solved")
	}
}

func TestMusicTriviaLettersRemainTypeable(t *testing.T) {
	s := &MusicTriviaSession{
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName: "Pulp",
				TrackName:  "Help the Aged",
				PreviewURL: "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessSong,
		},
	}

	for _, key := range []string{"h", "e", "l", "p"} {
		s.HandleKey(keyPress(key))
	}

	if s.GuessInput != "help" {
		t.Fatalf("expected h and p to be typeable, got %q", s.GuessInput)
	}
}

func TestMusicTriviaHintUnlocksAfterWrongGuesses(t *testing.T) {
	s := &MusicTriviaSession{
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName:   "Kate Bush",
				TrackName:    "Running Up That Hill",
				Genre:        "Pop",
				PrimaryGenre: "Alternative",
				PreviewURL:   "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessSong,
		},
	}

	s.revealHint()
	if s.HintCount != 0 {
		t.Fatal("expected hint to stay locked before guesses")
	}

	for range 2 {
		s.GuessInput = "wrong"
		s.submitGuess()
	}

	s.HandleKey(keyPress("?"))

	if s.HintCount != 1 {
		t.Fatalf("expected one hint, got %d", s.HintCount)
	}
}

func TestMusicTriviaEscapeClears(t *testing.T) {
	s := &MusicTriviaSession{
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName: "Prince",
				TrackName:  "Kiss",
				PreviewURL: "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessArtist,
		},
	}

	response := s.HandleKey(keyPress("esc"))

	if !response.Clear {
		t.Fatal("expected escape to clear the session")
	}

	if response.Message != "Music trivia cleared." {
		t.Fatalf("expected clear message, got %q", response.Message)
	}
}
