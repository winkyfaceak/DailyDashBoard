package app

import (
	"DailyDashBoard/internal/daily"
	"testing"
)

func TestMusicTriviaGuessFlow(t *testing.T) {
	m := InitialModel()
	m.musicTrivia = musicTriviaState{
		Active: true,
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName: "Prince",
				TrackName:  "Kiss",
				PreviewURL: "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessArtist,
		},
	}
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	for _, key := range []string{"P", "r", "i", "n", "c", "e"} {
		updated, _ := m.handleMusicTriviaKeyPress(keyPress(key))
		m = updated.(Model)
	}

	updated, _ := m.handleMusicTriviaKeyPress(keyPress("enter"))
	m = updated.(Model)

	if !m.musicTrivia.Solved {
		t.Fatal("expected music trivia to be solved")
	}
}

func TestMusicTriviaLettersRemainTypeable(t *testing.T) {
	m := InitialModel()
	m.musicTrivia = musicTriviaState{
		Active: true,
		Question: daily.MusicTriviaQuestion{
			Pick: daily.MusicPick{
				ArtistName: "Pulp",
				TrackName:  "Help the Aged",
				PreviewURL: "https://example.com/preview.m4a",
			},
			GuessKind: daily.MusicTriviaGuessSong,
		},
	}
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	for _, key := range []string{"h", "e", "l", "p"} {
		updated, _ := m.handleMusicTriviaKeyPress(keyPress(key))
		m = updated.(Model)
	}

	if m.musicTrivia.GuessInput != "help" {
		t.Fatalf("expected h and p to be typeable, got %q", m.musicTrivia.GuessInput)
	}
}

func TestMusicTriviaHintUnlocksAfterWrongGuesses(t *testing.T) {
	m := InitialModel()
	m.musicTrivia = musicTriviaState{
		Active: true,
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
	m.content.MusicTrivia = m.formatMusicTriviaContent()

	updated, _ := m.revealMusicTriviaHint()
	m = updated.(Model)
	if m.musicTrivia.HintCount != 0 {
		t.Fatal("expected hint to stay locked before guesses")
	}

	for range 2 {
		m.musicTrivia.GuessInput = "wrong"
		updated, _ = m.submitMusicTriviaGuess()
		m = updated.(Model)
	}

	updated, _ = m.handleMusicTriviaKeyPress(keyPress("?"))
	m = updated.(Model)

	if m.musicTrivia.HintCount != 1 {
		t.Fatalf("expected one hint, got %d", m.musicTrivia.HintCount)
	}
}
