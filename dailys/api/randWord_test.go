package api

import "testing"

func TestRandomWordEnglishUsesLocalWordList(t *testing.T) {
	word, err := RandomWord("en")
	if err != nil {
		t.Fatal(err)
	}

	if word == "" {
		t.Fatal("expected a random word")
	}
}

func TestIsEnglishLanguage(t *testing.T) {
	for _, language := range []string{"", "en", " EN "} {
		if !isEnglishLanguage(language) {
			t.Fatalf("expected %q to be treated as English", language)
		}
	}

	if isEnglishLanguage("es") {
		t.Fatal("expected es not to be treated as English")
	}
}
