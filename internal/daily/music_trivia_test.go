package daily

import "testing"

func TestMusicTriviaQuestionChecksArtistGuess(t *testing.T) {
	question := MusicTriviaQuestion{
		Pick: MusicPick{
			ArtistName:   "The Beatles",
			TrackName:    "Award Tour",
			Genre:        "Hip-Hop",
			PrimaryGenre: "Hip-Hop/Rap",
		},
		GuessKind: MusicTriviaGuessArtist,
	}

	if !question.CheckGuess("beatles") {
		t.Fatal("expected artist guess without leading article to be accepted")
	}

	if question.CheckGuess("De La Soul") {
		t.Fatal("expected wrong artist guess to be rejected")
	}
}

func TestMusicTriviaQuestionChecksSongGuess(t *testing.T) {
	question := MusicTriviaQuestion{
		Pick: MusicPick{
			ArtistName: "Kate Bush",
			TrackName:  "Running Up That Hill (A Deal With God)",
			AlbumName:  "Hounds of Love",
		},
		GuessKind: MusicTriviaGuessSong,
	}

	if !question.CheckGuess("running up that hill") {
		t.Fatal("expected song guess to be accepted")
	}

	if question.CheckGuess("wuthering heights") {
		t.Fatal("expected wrong song guess to be rejected")
	}
}

func TestMusicTriviaQuestionHints(t *testing.T) {
	question := MusicTriviaQuestion{
		Pick: MusicPick{
			ArtistName:   "Kate Bush",
			TrackName:    "Running Up That Hill",
			AlbumName:    "Hounds of Love",
			Genre:        "Pop",
			PrimaryGenre: "Alternative",
		},
		GuessKind: MusicTriviaGuessSong,
	}

	hints := question.Hints()
	if len(hints) != 3 {
		t.Fatalf("expected three hints, got %#v", hints)
	}

	if hints[0] != "Artist: Kate Bush" {
		t.Fatalf("expected artist hint, got %q", hints[0])
	}
}
