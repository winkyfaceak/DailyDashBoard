package daily

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	TarotOrientationUpright  = "upright"
	TarotOrientationReversed = "reversed"
)

var randomTarotCardURL = "https://freehoroscopeapi.com/api/v1/tarot/cards/random"

type TarotReading struct {
	Card        TarotCard
	Orientation string
	Meaning     string
}

type TarotCard struct {
	Name            string
	NameShort       string
	Value           string
	ValueInt        int
	Suit            string
	Type            string
	MeaningUp       string
	MeaningReversed string
	Description     string
}

type randomTarotCardResponse struct {
	Cards []tarotCardResponse `json:"cards"`
}

type tarotCardResponse struct {
	Name            string `json:"name"`
	NameShort       string `json:"name_short"`
	Value           string `json:"value"`
	ValueInt        int    `json:"value_int"`
	Suit            string `json:"suit"`
	Type            string `json:"type"`
	MeaningUp       string `json:"meaning_up"`
	MeaningReversed string `json:"meaning_rev"`
	Description     string `json:"desc"`
}

func RandomTarotReading() (TarotReading, error) {
	endpoint, err := randomTarotCardEndpoint()
	if err != nil {
		return TarotReading{}, fmt.Errorf("error parsing tarot card URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return TarotReading{}, fmt.Errorf("error creating tarot card request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return TarotReading{}, fmt.Errorf("error sending tarot card request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TarotReading{}, fmt.Errorf("tarot card API error: %s", resp.Status)
	}

	card, err := decodeRandomTarotCard(resp.Body)
	if err != nil {
		return TarotReading{}, err
	}

	orientation, err := randomTarotOrientation()
	if err != nil {
		return TarotReading{}, err
	}

	return tarotReadingFromCard(card, orientation), nil
}

func randomTarotCardEndpoint() (*url.URL, error) {
	endpoint, err := url.Parse(randomTarotCardURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("n", "1")
	query.Set("minor", "true")
	endpoint.RawQuery = query.Encode()

	return endpoint, nil
}

func decodeRandomTarotCard(reader io.Reader) (TarotCard, error) {
	var response randomTarotCardResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return TarotCard{}, fmt.Errorf("error decoding tarot card response: %w", err)
	}

	if len(response.Cards) == 0 {
		return TarotCard{}, fmt.Errorf("tarot card API returned no cards")
	}

	card := tarotCardFromResponse(response.Cards[0])
	if card.Name == "" {
		return TarotCard{}, fmt.Errorf("tarot card API returned no card name")
	}

	if card.MeaningUp == "" && card.MeaningReversed == "" {
		return TarotCard{}, fmt.Errorf("tarot card API returned no reading")
	}

	return card, nil
}

func tarotCardFromResponse(response tarotCardResponse) TarotCard {
	return TarotCard{
		Name:            cleanTarotText(response.Name),
		NameShort:       cleanTarotText(response.NameShort),
		Value:           cleanTarotText(response.Value),
		ValueInt:        response.ValueInt,
		Suit:            cleanTarotText(response.Suit),
		Type:            cleanTarotText(response.Type),
		MeaningUp:       cleanTarotText(response.MeaningUp),
		MeaningReversed: cleanTarotText(response.MeaningReversed),
		Description:     cleanTarotText(response.Description),
	}
}

func randomTarotOrientation() (string, error) {
	index, err := randomInt(2)
	if err != nil {
		return "", err
	}

	if index == 1 {
		return TarotOrientationReversed, nil
	}

	return TarotOrientationUpright, nil
}

func tarotReadingFromCard(card TarotCard, orientation string) TarotReading {
	if orientation == "" {
		orientation = TarotOrientationUpright
	}

	return TarotReading{
		Card:        card,
		Orientation: orientation,
		Meaning:     tarotMeaning(card, orientation),
	}
}

func tarotMeaning(card TarotCard, orientation string) string {
	if orientation == TarotOrientationReversed && card.MeaningReversed != "" {
		return card.MeaningReversed
	}

	if card.MeaningUp != "" {
		return card.MeaningUp
	}

	return card.MeaningReversed
}

func cleanTarotText(value string) string {
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}
