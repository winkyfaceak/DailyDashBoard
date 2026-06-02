package daily

import (
	"testing"
)

func TestRandomMusicFromGenresRequiresGenres(t *testing.T) {
	_, err := RandomMusicFromGenres(nil)
	if err == nil {
		t.Fatal("expected an error for an empty genre list")
	}
}

func TestHighResolutionArtworkURL(t *testing.T) {
	imageURL := "https://is1-ssl.mzstatic.com/image/thumb/Music124/v4/example/cover.jpg/100x100bb.jpg"
	want := "https://is1-ssl.mzstatic.com/image/thumb/Music124/v4/example/cover.jpg/600x600bb.jpg"

	if got := highResolutionArtworkURL(imageURL); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestMusicPickFromTrackCleansMetadata(t *testing.T) {
	pick := musicPickFromTrack(MusicGenre{Name: "Jazz"}, itunesTrack{
		Kind:         "song",
		TrackName:    " Take Five ",
		ArtistName:   " Dave Brubeck ",
		Collection:   " Time Out ",
		PrimaryGenre: "Jazz",
		ReleaseDate:  "1959-12-14T12:00:00Z",
		PreviewURL:   " https://example.com/preview.m4a ",
		ArtworkURL:   "https://example.com/100x100bb.jpg",
		TrackURL:     " https://example.com/song ",
	})

	if pick.TrackName != "Take Five" {
		t.Fatalf("expected cleaned track name, got %q", pick.TrackName)
	}

	if pick.ArtistName != "Dave Brubeck" {
		t.Fatalf("expected cleaned artist name, got %q", pick.ArtistName)
	}

	if pick.ReleaseDate.Year() != 1959 {
		t.Fatalf("expected release year 1959, got %d", pick.ReleaseDate.Year())
	}

	if pick.PreviewURL != "https://example.com/preview.m4a" {
		t.Fatalf("expected cleaned preview URL, got %q", pick.PreviewURL)
	}
}
