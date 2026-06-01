package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	randomWordURL        = "https://random-word-api.herokuapp.com/word"
	randomWordAPITimeout = 2 * time.Second
)

var englishRandomWords = []string{
	"abundance",
	"adventure",
	"anchor",
	"aperture",
	"blossom",
	"cadence",
	"catalyst",
	"clarity",
	"compass",
	"constellation",
	"curious",
	"diligent",
	"dusk",
	"echo",
	"eloquent",
	"ember",
	"empathy",
	"evergreen",
	"flourish",
	"gentle",
	"glimmer",
	"gratitude",
	"harbor",
	"harmony",
	"horizon",
	"insight",
	"intricate",
	"jubilant",
	"lantern",
	"legacy",
	"luminous",
	"meadow",
	"momentum",
	"mosaic",
	"nimble",
	"nuance",
	"orchard",
	"patience",
	"pebble",
	"perspective",
	"practical",
	"quiet",
	"radiant",
	"resilient",
	"river",
	"sanctuary",
	"serendipity",
	"sincere",
	"solitude",
	"spark",
	"steady",
	"tapestry",
	"thoughtful",
	"threshold",
	"tranquil",
	"twilight",
	"vibrant",
	"voyage",
	"wonder",
	"zephyr",
}

func RandomWord(language string) (string, error) {
	if isEnglishLanguage(language) {
		return randomEnglishWord()
	}

	return randomWordFromAPI(language)
}

func randomWordFromAPI(language string) (string, error) {
	endpoint, err := url.Parse(randomWordURL)
	if err != nil {
		return "", fmt.Errorf("error parsing random word URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("number", "1")

	if language != "" && language != "en" {
		query.Set("lang", language)
	}

	endpoint.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), randomWordAPITimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating random word request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending random word request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("random word API error: %s", resp.Status)
	}

	var words []string

	err = json.NewDecoder(resp.Body).Decode(&words)
	if err != nil {
		return "", fmt.Errorf("error decoding random word response: %w", err)
	}

	if len(words) == 0 {
		return "", fmt.Errorf("random word API returned no words")
	}

	word := strings.TrimSpace(words[0])
	if word == "" {
		return "", fmt.Errorf("random word API returned an empty word")
	}

	return word, nil
}

func isEnglishLanguage(language string) bool {
	language = strings.TrimSpace(strings.ToLower(language))

	return language == "" || language == "en"
}

func randomEnglishWord() (string, error) {
	index, err := randomInt(len(englishRandomWords))
	if err != nil {
		return "", err
	}

	return englishRandomWords[index], nil
}

func RandomDictionarySet(language string) ([]DictionaryEntry, error) {
	const maxAttempts = 5

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		word, err := RandomWord(language)
		if err != nil {
			return nil, err
		}

		entries, err := DictionarySet(language, word)
		if err == nil {
			return entries, nil
		}

		lastErr = fmt.Errorf("%s: %w", word, err)
	}

	return nil, fmt.Errorf("could not find a dictionary entry after %d random words: %w", maxAttempts, lastErr)
}
