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
	tarotDailyCardCount      = 3
)

var randomTarotCardURL = "https://freehoroscopeapi.com/api/v1/tarot/cards/random"

type TarotReading struct {
	Card        TarotCard
	Orientation string
	Meaning     string
	SpreadName  string
	Cards       []TarotSpreadCard
	Synthesis   string
}

type TarotSpreadCard struct {
	Position      string
	PositionFocus string
	Card          TarotCard
	Orientation   string
	Meaning       string
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
	ImageURL        string
}

type tarotSpreadPosition struct {
	Name  string
	Focus string
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

var dailyTarotSpreadPositions = []tarotSpreadPosition{
	{
		Name:  "Past",
		Focus: "Roots, prior choices, habits, or events shaping today's situation.",
	},
	{
		Name:  "Present",
		Focus: "The current energy, challenge, or opportunity that needs honest attention.",
	},
	{
		Name:  "Future",
		Focus: "The likely direction or advice if the present pattern continues.",
	},
}

var tarotMajorImageFiles = map[int]string{
	0:  "RWS_Tarot_00_Fool.jpg",
	1:  "RWS_Tarot_01_Magician.jpg",
	2:  "RWS_Tarot_02_High_Priestess.jpg",
	3:  "RWS_Tarot_03_Empress.jpg",
	4:  "RWS_Tarot_04_Emperor.jpg",
	5:  "RWS_Tarot_05_Hierophant.jpg",
	6:  "RWS_Tarot_06_Lovers.jpg",
	7:  "RWS_Tarot_07_Chariot.jpg",
	8:  "RWS_Tarot_08_Strength.jpg",
	9:  "RWS_Tarot_09_Hermit.jpg",
	10: "RWS_Tarot_10_Wheel_of_Fortune.jpg",
	11: "RWS_Tarot_11_Justice.jpg",
	12: "RWS_Tarot_12_Hanged_Man.jpg",
	13: "RWS_Tarot_13_Death.jpg",
	14: "RWS_Tarot_14_Temperance.jpg",
	15: "RWS_Tarot_15_Devil.jpg",
	16: "RWS_Tarot_16_Tower.jpg",
	17: "RWS_Tarot_17_Star.jpg",
	18: "RWS_Tarot_18_Moon.jpg",
	19: "RWS_Tarot_19_Sun.jpg",
	20: "RWS_Tarot_20_Judgement.jpg",
	21: "RWS_Tarot_21_World.jpg",
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

	cards, err := decodeRandomTarotCards(resp.Body)
	if err != nil {
		return TarotReading{}, err
	}

	if len(cards) < tarotDailyCardCount {
		return TarotReading{}, fmt.Errorf("tarot card API returned %d cards; need %d", len(cards), tarotDailyCardCount)
	}

	spreadCards := make([]TarotSpreadCard, 0, tarotDailyCardCount)
	for index, card := range cards[:tarotDailyCardCount] {
		orientation, err := randomTarotOrientation()
		if err != nil {
			return TarotReading{}, err
		}

		position := dailyTarotSpreadPositions[index]
		spreadCards = append(spreadCards, tarotSpreadCardFromCard(card, orientation, position))
	}

	return tarotReadingFromSpreadCards("Past / Present / Future", spreadCards), nil
}

func randomTarotCardEndpoint() (*url.URL, error) {
	endpoint, err := url.Parse(randomTarotCardURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("n", "3")
	query.Set("minor", "true")
	endpoint.RawQuery = query.Encode()

	return endpoint, nil
}

func decodeRandomTarotCard(reader io.Reader) (TarotCard, error) {
	cards, err := decodeRandomTarotCards(reader)
	if err != nil {
		return TarotCard{}, err
	}

	return cards[0], nil
}

func decodeRandomTarotCards(reader io.Reader) ([]TarotCard, error) {
	var response randomTarotCardResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding tarot card response: %w", err)
	}

	if len(response.Cards) == 0 {
		return nil, fmt.Errorf("tarot card API returned no cards")
	}

	cards := make([]TarotCard, 0, len(response.Cards))
	for _, item := range response.Cards {
		card := tarotCardFromResponse(item)
		if card.Name == "" {
			return nil, fmt.Errorf("tarot card API returned no card name")
		}

		if card.MeaningUp == "" && card.MeaningReversed == "" {
			return nil, fmt.Errorf("tarot card API returned no reading")
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func tarotCardFromResponse(response tarotCardResponse) TarotCard {
	card := TarotCard{
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
	card.ImageURL = TarotCardImageURL(card)

	return card
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

	spreadCard := tarotSpreadCardFromCard(card, orientation, tarotSpreadPosition{
		Name:  "Card",
		Focus: "Daily guidance for the question or day ahead.",
	})

	return tarotReadingFromSpreadCards("Single Card", []TarotSpreadCard{spreadCard})
}

func tarotSpreadCardFromCard(card TarotCard, orientation string, position tarotSpreadPosition) TarotSpreadCard {
	if orientation == "" {
		orientation = TarotOrientationUpright
	}

	return TarotSpreadCard{
		Position:      position.Name,
		PositionFocus: position.Focus,
		Card:          card,
		Orientation:   orientation,
		Meaning:       tarotMeaning(card, orientation),
	}
}

func tarotReadingFromSpreadCards(spreadName string, cards []TarotSpreadCard) TarotReading {
	reading := TarotReading{
		SpreadName: spreadName,
		Cards:      cards,
		Synthesis:  tarotSynthesis(cards),
	}

	if len(cards) > 0 {
		reading.Card = cards[0].Card
		reading.Orientation = cards[0].Orientation
		reading.Meaning = cards[0].Meaning
	}

	return reading
}

func tarotSynthesis(cards []TarotSpreadCard) string {
	if len(cards) == 0 {
		return ""
	}

	if len(cards) == 1 {
		return fmt.Sprintf("%s offers the central guidance for the day. Read the card as a prompt for reflection and action rather than a fixed prediction.", tarotSpreadCardLabel(cards[0]))
	}

	if len(cards) >= tarotDailyCardCount {
		return fmt.Sprintf(
			"The story moves from %s into %s, with %s as the next opening. Give the present card the most weight, and treat the future card as trajectory or advice rather than a fixed prediction.",
			tarotSpreadCardLabel(cards[0]),
			tarotSpreadCardLabel(cards[1]),
			tarotSpreadCardLabel(cards[2]),
		)
	}

	return "Read the cards in their chosen positions, then connect them into one story."
}

func tarotSpreadCardLabel(card TarotSpreadCard) string {
	label := card.Card.Name
	if card.Orientation != "" {
		label += " " + titleCaseTarot(card.Orientation)
	}

	return label
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

func TarotCardImageURL(card TarotCard) string {
	fileName := tarotCardImageFileName(card)
	if fileName == "" {
		return ""
	}

	return "https://commons.wikimedia.org/wiki/Special:FilePath/" + url.PathEscape(fileName) + "?width=360"
}

func tarotCardImageFileName(card TarotCard) string {
	if strings.EqualFold(card.Type, "major") {
		return tarotMajorImageFiles[card.ValueInt]
	}

	prefix := tarotSuitImagePrefix(card.Suit)
	if prefix == "" {
		return ""
	}

	valueInt := card.ValueInt
	if valueInt == 0 {
		valueInt = tarotValueInt(card.Value)
	}

	if valueInt < 1 || valueInt > 14 {
		return ""
	}

	if prefix == "Wands" && valueInt == 9 {
		return "Tarot_Nine_of_Wands.jpg"
	}

	return fmt.Sprintf("%s%02d.jpg", prefix, valueInt)
}

func tarotSuitImagePrefix(suit string) string {
	switch strings.ToLower(strings.TrimSpace(suit)) {
	case "wands", "wand":
		return "Wands"
	case "cups", "cup":
		return "Cups"
	case "swords", "sword":
		return "Swords"
	case "pentacles", "pentacle", "coins", "coin", "pents":
		return "Pents"
	default:
		return ""
	}
}

func tarotValueInt(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "ace":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	case "ten":
		return 10
	case "page":
		return 11
	case "knight":
		return 12
	case "queen":
		return 13
	case "king":
		return 14
	default:
		return 0
	}
}

func cleanTarotText(value string) string {
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}

func titleCaseTarot(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return strings.ToUpper(value[:1]) + strings.ToLower(value[1:])
}
