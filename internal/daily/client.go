// Package daily contains the API clients that supply the dashboard's
// content: facts, words, trivia, tarot, artwork, news, music, and recipes.
package daily

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiUserAgent = "DailyDashBoard/1.0"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// requestOption customizes an API request before it is sent.
type requestOption func(*http.Request)

func withHeader(name string, value string) requestOption {
	return func(req *http.Request) {
		req.Header.Set(name, value)
	}
}

// getJSON performs a GET against endpoint and hands the response body to
// decode. label names the API in error messages, e.g. "trivia" produces
// "error sending trivia request" and "trivia API error: ...".
func getJSON(endpoint string, label string, decode func(io.Reader) error, options ...requestOption) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating %s request: %w", label, err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	for _, option := range options {
		option(req)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending %s request: %w", label, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s API error: %s", label, resp.Status)
	}

	return decode(resp.Body)
}
