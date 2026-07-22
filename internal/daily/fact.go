package daily

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/url"
	"strings"
)

var randomFactURL = "https://uselessfacts.jsph.pl/api/v2/facts/random"

type Fact struct {
	Text      string
	Source    string
	SourceURL string
	Permalink string
	Language  string
}

type randomFactResponse struct {
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url"`
	Permalink string `json:"permalink"`
	Language  string `json:"language"`
}

func RandomFact(language string) (Fact, error) {
	endpoint, err := randomFactEndpoint(language)
	if err != nil {
		return Fact{}, fmt.Errorf("error parsing random fact URL: %w", err)
	}

	var fact Fact

	err = getJSON(endpoint.String(), "random fact", func(body io.Reader) error {
		decoded, decodeErr := decodeRandomFact(body)
		if decodeErr != nil {
			return decodeErr
		}

		fact = decoded

		return nil
	})
	if err != nil {
		return Fact{}, err
	}

	return fact, nil
}

func randomFactEndpoint(language string) (*url.URL, error) {
	endpoint, err := url.Parse(randomFactURL)
	if err != nil {
		return nil, err
	}

	language = strings.TrimSpace(strings.ToLower(language))
	if language == "" {
		language = "en"
	}

	query := endpoint.Query()
	query.Set("language", language)
	endpoint.RawQuery = query.Encode()

	return endpoint, nil
}

func decodeRandomFact(reader io.Reader) (Fact, error) {
	var response randomFactResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return Fact{}, fmt.Errorf("error decoding random fact response: %w", err)
	}

	fact := Fact{
		Text:      cleanFactText(response.Text),
		Source:    strings.TrimSpace(response.Source),
		SourceURL: strings.TrimSpace(response.SourceURL),
		Permalink: strings.TrimSpace(response.Permalink),
		Language:  strings.TrimSpace(response.Language),
	}

	if fact.Text == "" {
		return Fact{}, fmt.Errorf("random fact API returned no fact")
	}

	return fact, nil
}

func cleanFactText(value string) string {
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}
