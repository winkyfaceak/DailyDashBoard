package daily

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomFactUsesLanguageQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("language"); got != "en" {
			t.Fatalf("expected language query to be en, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected JSON accept header, got %q", got)
		}

		fmt.Fprint(w, `{"text":"The Eiffel Tower can grow taller in warm weather.","source":"uselessfacts","source_url":"https://example.com/source","permalink":"https://example.com/fact","language":"en"}`)
	}))
	defer server.Close()

	restoreRandomFactURL(t, server.URL)

	fact, err := RandomFact(" EN ")
	if err != nil {
		t.Fatal(err)
	}

	if fact.Text != "The Eiffel Tower can grow taller in warm weather." {
		t.Fatalf("expected decoded fact text, got %q", fact.Text)
	}

	if fact.SourceURL != "https://example.com/source" {
		t.Fatalf("expected source URL, got %q", fact.SourceURL)
	}
}

func TestDecodeRandomFactCleansText(t *testing.T) {
	fact, err := decodeRandomFact(strings.NewReader(`{"text":"  Bananas&nbsp;are berries.\n\n","language":"en"}`))
	if err != nil {
		t.Fatal(err)
	}

	if fact.Text != "Bananas are berries." {
		t.Fatalf("expected cleaned fact text, got %q", fact.Text)
	}
}

func TestDecodeRandomFactRejectsEmptyText(t *testing.T) {
	_, err := decodeRandomFact(strings.NewReader(`{"text":"   "}`))
	if err == nil {
		t.Fatal("expected empty fact text error")
	}
}

func restoreRandomFactURL(t *testing.T, testURL string) {
	t.Helper()

	originalURL := randomFactURL
	randomFactURL = testURL

	t.Cleanup(func() {
		randomFactURL = originalURL
	})
}
