package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomWordEnglishUsesAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected English request to omit query params, got %q", r.URL.RawQuery)
		}

		fmt.Fprint(w, `{"word":"zeppelin"}`)
	}))
	defer server.Close()

	restoreRandomWordURL(t, server.URL)

	word, err := RandomWord("en")
	if err != nil {
		t.Fatal(err)
	}

	if word != "zeppelin" {
		t.Fatalf("expected API word, got %q", word)
	}
}

func TestDecodeRandomWordSupportsStringArrays(t *testing.T) {
	word, err := decodeRandomWord(strings.NewReader(`["zeppelin"]`))
	if err != nil {
		t.Fatal(err)
	}

	if word != "zeppelin" {
		t.Fatalf("expected decoded word, got %q", word)
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

func restoreRandomWordURL(t *testing.T, testURL string) {
	t.Helper()

	originalURL := randomWordURL
	randomWordURL = testURL

	t.Cleanup(func() {
		randomWordURL = originalURL
	})
}
