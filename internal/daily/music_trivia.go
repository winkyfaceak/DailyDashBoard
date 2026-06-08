package daily

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	MusicTriviaGuessArtist = "artist"
	MusicTriviaGuessSong   = "song"
)

type MusicTriviaQuestion struct {
	Pick      MusicPick
	GuessKind string
}

func RandomMusicTriviaQuestion() (MusicTriviaQuestion, error) {
	const maxAttempts = 5

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		pick, err := RandomMusic()
		if err != nil {
			return MusicTriviaQuestion{}, err
		}

		if pick.PreviewURL == "" {
			lastErr = fmt.Errorf("%s by %s does not include a preview", pick.TrackName, pick.ArtistName)
			continue
		}

		guessKind, err := randomMusicTriviaGuessKind()
		if err != nil {
			return MusicTriviaQuestion{}, err
		}

		return MusicTriviaQuestion{
			Pick:      pick,
			GuessKind: guessKind,
		}, nil
	}

	return MusicTriviaQuestion{}, fmt.Errorf("could not find a music trivia preview: %w", lastErr)
}

func randomMusicTriviaGuessKind() (string, error) {
	index, err := randomInt(2)
	if err != nil {
		return "", err
	}

	if index == 1 {
		return MusicTriviaGuessSong, nil
	}

	return MusicTriviaGuessArtist, nil
}

func (q MusicTriviaQuestion) Prompt() string {
	if q.GuessKind == MusicTriviaGuessSong {
		return "What song is this?"
	}

	return "Who is the artist?"
}

func (q MusicTriviaQuestion) TargetAnswer() string {
	if q.GuessKind == MusicTriviaGuessSong {
		return q.Pick.TrackName
	}

	return q.Pick.ArtistName
}

func (q MusicTriviaQuestion) CheckGuess(guess string) bool {
	normalizedGuess := normalizeGuessText(guess)
	normalizedAnswer := normalizeGuessText(q.TargetAnswer())
	if normalizedGuess == "" || normalizedAnswer == "" {
		return false
	}

	if normalizedGuess == normalizedAnswer {
		return true
	}

	if strings.TrimPrefix(normalizedAnswer, "the ") == normalizedGuess || strings.TrimPrefix(normalizedAnswer, "a ") == normalizedGuess {
		return true
	}

	return len(normalizedGuess) >= 4 && strings.Contains(normalizedAnswer, normalizedGuess)
}

func (q MusicTriviaQuestion) Hints() []string {
	answer := q.TargetAnswer()
	firstLetter := firstLetter(answer)
	length := visibleLetterCount(answer)

	if q.GuessKind == MusicTriviaGuessSong {
		hints := []string{
			fmt.Sprintf("Artist: %s", q.Pick.ArtistName),
			fmt.Sprintf("The song title starts with %q and has %d letters or numbers.", firstLetter, length),
		}

		if q.Pick.AlbumName != "" {
			hints = append(hints, fmt.Sprintf("Album: %s", q.Pick.AlbumName))
		}

		return hints
	}

	hints := []string{
		fmt.Sprintf("Genre: %s", fallbackString(q.Pick.PrimaryGenre, q.Pick.Genre)),
		fmt.Sprintf("The artist starts with %q and has %d letters or numbers.", firstLetter, length),
		fmt.Sprintf("Song title: %s", q.Pick.TrackName),
	}

	return hints
}

func normalizeGuessText(value string) string {
	var builder strings.Builder
	previousWasSpace := true

	for _, char := range strings.ToLower(value) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(char)
			previousWasSpace = false
			continue
		}

		if !previousWasSpace {
			builder.WriteRune(' ')
			previousWasSpace = true
		}
	}

	return strings.TrimSpace(builder.String())
}

func firstLetter(value string) string {
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			return strings.ToUpper(string(char))
		}
	}

	return "?"
}

func visibleLetterCount(value string) int {
	count := 0
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}

	return count
}
