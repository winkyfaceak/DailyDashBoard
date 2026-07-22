package daily

import (
	"encoding/json"
	"fmt"
	"io"
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

	var entries []DictionaryEntry

	err := getJSON(endpoint, "dictionary", func(body io.Reader) error {
		if err := json.NewDecoder(body).Decode(&entries); err != nil {
			return fmt.Errorf("error decoding dictionary response: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}
