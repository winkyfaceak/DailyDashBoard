package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	legacyRandomWordURL  = "https://random-word-api.herokuapp.com/word"
	randomWordAPITimeout = 10 * time.Second
)

var randomWordURL = "https://www.wordgamedb.com/api/v2/words/random"

func RandomWord(language string) (string, error) {
	return randomWordFromAPI(language)
}

func randomWordFromAPI(language string) (string, error) {
	language = strings.TrimSpace(strings.ToLower(language))

	endpoint, err := randomWordEndpoint(language)
	if err != nil {
		return "", fmt.Errorf("error parsing random word URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), randomWordAPITimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating random word request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending random word request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("random word API error: %s", resp.Status)
	}

	word, err := decodeRandomWord(resp.Body)
	if err != nil {
		return "", err
	}

	return word, nil
}

func randomWordEndpoint(language string) (*url.URL, error) {
	if isEnglishLanguage(language) {
		return url.Parse(randomWordURL)
	}

	endpoint, err := url.Parse(legacyRandomWordURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("number", "1")
	query.Set("lang", language)
	endpoint.RawQuery = query.Encode()

	return endpoint, nil
}

func decodeRandomWord(reader io.Reader) (string, error) {
	var payload json.RawMessage
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return "", fmt.Errorf("error decoding random word response: %w", err)
	}

	var object struct {
		Word string `json:"word"`
	}
	if err := json.Unmarshal(payload, &object); err == nil && strings.TrimSpace(object.Word) != "" {
		return strings.TrimSpace(object.Word), nil
	}

	var words []string
	if err := json.Unmarshal(payload, &words); err == nil && len(words) > 0 && strings.TrimSpace(words[0]) != "" {
		return strings.TrimSpace(words[0]), nil
	}

	return "", fmt.Errorf("random word API returned no words")
}

func isEnglishLanguage(language string) bool {
	language = strings.TrimSpace(strings.ToLower(language))

	return language == "" || language == "en"
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
