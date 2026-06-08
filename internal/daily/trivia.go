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

var randomTriviaURL = "https://opentdb.com/api.php"

type TriviaQuestion struct {
	Category           string
	Difficulty         string
	Question           string
	CorrectAnswer      string
	Choices            []string
	CorrectChoiceIndex int
}

type randomTriviaResponse struct {
	ResponseCode int              `json:"response_code"`
	Results      []triviaResponse `json:"results"`
}

type triviaResponse struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func RandomTriviaQuestion() (TriviaQuestion, error) {
	endpoint, err := randomTriviaEndpoint()
	if err != nil {
		return TriviaQuestion{}, fmt.Errorf("error parsing trivia URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return TriviaQuestion{}, fmt.Errorf("error creating trivia request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return TriviaQuestion{}, fmt.Errorf("error sending trivia request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TriviaQuestion{}, fmt.Errorf("trivia API error: %s", resp.Status)
	}

	question, err := decodeRandomTriviaQuestion(resp.Body)
	if err != nil {
		return TriviaQuestion{}, err
	}

	return question, nil
}

func randomTriviaEndpoint() (*url.URL, error) {
	endpoint, err := url.Parse(randomTriviaURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("amount", "1")
	query.Set("type", "multiple")
	endpoint.RawQuery = query.Encode()

	return endpoint, nil
}

func decodeRandomTriviaQuestion(reader io.Reader) (TriviaQuestion, error) {
	var response randomTriviaResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return TriviaQuestion{}, fmt.Errorf("error decoding trivia response: %w", err)
	}

	if response.ResponseCode != 0 {
		return TriviaQuestion{}, fmt.Errorf("trivia API response code: %d", response.ResponseCode)
	}

	if len(response.Results) == 0 {
		return TriviaQuestion{}, fmt.Errorf("trivia API returned no questions")
	}

	question, err := triviaQuestionFromResponse(response.Results[0])
	if err != nil {
		return TriviaQuestion{}, err
	}

	return question, nil
}

func triviaQuestionFromResponse(response triviaResponse) (TriviaQuestion, error) {
	correctAnswer := cleanTriviaText(response.CorrectAnswer)
	if correctAnswer == "" {
		return TriviaQuestion{}, fmt.Errorf("trivia API returned no correct answer")
	}

	questionText := cleanTriviaText(response.Question)
	if questionText == "" {
		return TriviaQuestion{}, fmt.Errorf("trivia API returned no question")
	}

	choices := make([]string, 0, len(response.IncorrectAnswers)+1)
	choices = append(choices, correctAnswer)
	for _, answer := range response.IncorrectAnswers {
		cleanedAnswer := cleanTriviaText(answer)
		if cleanedAnswer != "" {
			choices = append(choices, cleanedAnswer)
		}
	}

	if len(choices) < 2 {
		return TriviaQuestion{}, fmt.Errorf("trivia API returned too few choices")
	}

	correctChoiceIndex, err := shuffleTriviaChoices(choices, correctAnswer)
	if err != nil {
		return TriviaQuestion{}, err
	}

	return TriviaQuestion{
		Category:           cleanTriviaText(response.Category),
		Difficulty:         cleanTriviaText(response.Difficulty),
		Question:           questionText,
		CorrectAnswer:      correctAnswer,
		Choices:            choices,
		CorrectChoiceIndex: correctChoiceIndex,
	}, nil
}

func shuffleTriviaChoices(choices []string, correctAnswer string) (int, error) {
	correctChoiceIndex := 0
	for index := len(choices) - 1; index > 0; index-- {
		swapIndex, err := randomInt(index + 1)
		if err != nil {
			return 0, err
		}

		choices[index], choices[swapIndex] = choices[swapIndex], choices[index]
	}

	for index, choice := range choices {
		if choice == correctAnswer {
			correctChoiceIndex = index
			break
		}
	}

	return correctChoiceIndex, nil
}

func (q TriviaQuestion) CheckChoice(index int) bool {
	return index == q.CorrectChoiceIndex
}

func cleanTriviaText(value string) string {
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}
