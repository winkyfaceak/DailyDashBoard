package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DictionaryEntry struct {
	Word      string     `json:"word"`
	Phonetic  string     `json:"phonetic"`
	Phonetics []Phonetic `json:"phonetics"`
	Origin    string     `json:"origin"`
	Meanings  []Meaning  `json:"meanings"`
}

type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}

func DictionarySet(language string, word string) ([]DictionaryEntry, error) {
	endpoint := fmt.Sprintf(
		"https://api.dictionaryapi.dev/api/v2/entries/%s/%s",
		url.PathEscape(language),
		url.PathEscape(word),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating dictionary request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending dictionary request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dictionary API error: %s", resp.Status)
	}

	var entries []DictionaryEntry

	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("error decoding dictionary response: %w", err)
	}

	return entries, nil
}
